package api

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/xh3sh/go-real-time-chats/internal/emmit"
	"github.com/xh3sh/go-real-time-chats/internal/templates"
)

type sseHandler struct {
	emitter *emmit.Emitter
}

func NewSSEHandler(emmiter *emmit.Emitter) *sseHandler {
	return &sseHandler{emitter: emmiter}
}

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
	IsSelf   bool   `json:"isSelf"`
}

type SSEClient struct {
	Username string
	Channel  chan emmit.Message
}

var (
	sseClients    = make(map[*SSEClient]struct{})
	sseClientsMux sync.Mutex
)

func (s *sseHandler) SSEHandler(c echo.Context) error {
	username := c.QueryParam("username")
	if username == "" {
		username = "Anonymous"
	}

	client := &SSEClient{
		Username: username,
		Channel:  make(chan emmit.Message),
	}

	sseClientsMux.Lock()
	sseClients[client] = struct{}{}
	sseClientsMux.Unlock()

	defer func() {
		sseClientsMux.Lock()
		delete(sseClients, client)
		sseClientsMux.Unlock()
		close(client.Channel)
	}()

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")

	tmpl := templates.GetTemplates(c)

	off := s.emitter.On("newMessage", func(msg emmit.Message) {
		msg.IsSelf = msg.Username == client.Username

		var htmlBuffer bytes.Buffer
		if err := tmpl.ExecuteTemplate(&htmlBuffer, "user-message", msg); err != nil {
			log.Printf("failed to render template: %v\n", err)
			return
		}

		htmlContent := htmlBuffer.String()
		htmlContent = strings.ReplaceAll(htmlContent, "\r\n", "\n") // Windows -> Unix
		htmlContent = strings.ReplaceAll(htmlContent, "\n", "")     // Убираем все переводы строки
		if _, err := c.Response().Write([]byte("data:" + htmlContent + "\n\n")); err != nil {
			log.Println(err)
		}
		c.Response().Flush()
	})

	defer off()

	<-c.Request().Context().Done()
	return nil
}

func (s *sseHandler) SSEMessageHandler(c echo.Context) error {
	message := c.FormValue("content")
	username := c.FormValue("username")
	if message == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":   "Empty message",
			"message": "Message content cannot be empty.",
		})
	}

	msg := emmit.Message{
		Username: username,
		Content:  message,
	}

	s.emitter.Emit("newMessage", msg)

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Message broadcasted successfully.",
	})
}
