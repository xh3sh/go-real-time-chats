package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xh3sh/go-real-time-chats/internal/emmit"
)

type longPollHandler struct {
	emitter *emmit.Emitter
}

func NewLongPollHandler(emmiter *emmit.Emitter) *longPollHandler {
	return &longPollHandler{emitter: emmiter}
}

func (l *longPollHandler) GetMessages(c echo.Context) error {
	currentUser := c.FormValue("username")
	msg, ok := l.emitter.Wait("newMessage", 60*time.Second)
	if ok {
		msg.IsSelf = msg.Username == currentUser

		return c.Render(http.StatusOK, "user-message", msg)
	}

	return c.NoContent(http.StatusNoContent)
}

func (l *longPollHandler) PostMessage(c echo.Context) error {
	message := c.FormValue("message")
	username := c.FormValue("username")

	if username == "" || message == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username and message cannot be empty"})
	}

	msg := emmit.Message{
		Username: username,
		Content:  message,
	}

	l.emitter.Emit("newMessage", msg)

	return c.JSON(http.StatusOK, map[string]string{"message": message})
}
