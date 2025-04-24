package api

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/xh3sh/go-real-time-chats/internal/emmit"
)

type websocketHandler struct {
	emitter *emmit.Emitter
}

func NewWebsocketHandler(emmiter *emmit.Emitter) *websocketHandler {
	return &websocketHandler{emitter: emmiter}
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

	off := w.emitter.On("newMessage", func(msg emmit.Message) {
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

		htmlMessage := htmlBuffer.Bytes()
		if err := client.Conn.WriteMessage(websocket.TextMessage, htmlMessage); err != nil {
			client.Conn.Close()
			clientsMux.Lock()
			delete(clients, client.Conn)
			clientsMux.Unlock()
		}
	})

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
		w.emitter.Emit("newMessage", msg)
	}

	return nil
}
