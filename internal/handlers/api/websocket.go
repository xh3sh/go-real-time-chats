package api

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/xh3sh/go-real-time-chats/internal/emmit"
	"github.com/xh3sh/go-real-time-chats/internal/repo"
)

type websocketHandler struct {
	repo    *repo.RedisRepository
	emitter *emmit.Emitter
}

func NewWebsocketHandler(repo *repo.RedisRepository, emmiter *emmit.Emitter) *websocketHandler {
	return &websocketHandler{repo: repo, emitter: emmiter}
}

type Client struct {
	Conn     *websocket.Conn
	Username string
}

var (
	upgrader   = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	clients    = make(map[*websocket.Conn]*Client)
	clientsMux sync.Mutex
)

func (w *websocketHandler) WebSocketHandler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	username := c.QueryParam("username")
	if username == "" {
		username = "Anonymous"
	}

	client := &Client{Conn: conn, Username: username}

	clientsMux.Lock()
	clients[conn] = client
	clientsMux.Unlock()

	sendMsg := func(msg emmit.Message) {
		msg.IsSelf = msg.Username == client.Username

		tmpl, err := template.ParseFiles("web/views/components/user-message.html")
		if err != nil {
			log.Printf("template parsing error: %s", err)
			return
		}

		var htmlBuffer bytes.Buffer
		if err := tmpl.ExecuteTemplate(&htmlBuffer, "user-message", msg); err != nil {
			log.Printf("failed to render template: %v\n", err)
			return
		}

		htmlContent := htmlBuffer.String()
		htmlContent = strings.Replace(htmlContent, "class=\"message-wrapper\"", "hx-swap-oob=\"beforeend:#chat-message\" class=\"message-wrapper\"", 1)

		if err := client.Conn.WriteMessage(websocket.TextMessage, []byte(htmlContent)); err != nil {
			client.Conn.Close()
			clientsMux.Lock()
			delete(clients, client.Conn)
			clientsMux.Unlock()
		}
	}

	if messages, err := w.repo.GetMessages(c.Request().Context()); err == nil {
		for _, msg := range messages {
			sendMsg(msg)
		}
	}

	off := w.emitter.On("newMessage", sendMsg)

	defer off()

	for {
		var msg emmit.Message
		if err := conn.ReadJSON(&msg); err != nil {
			clientsMux.Lock()
			delete(clients, conn)
			clientsMux.Unlock()
			break
		}
		msg.Username = client.Username
		msg.Timestamp = time.Now().Unix()

		if err := w.repo.SaveMessage(c.Request().Context(), msg); err != nil {
			log.Printf("failed to save message to redis: %v", err)
		}

		w.emitter.Emit("newMessage", msg)
	}

	return nil
}
