package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xh3sh/go-real-time-chats/internal/emmit"
	"github.com/xh3sh/go-real-time-chats/internal/repo"
)

type longPollHandler struct {
	repo    *repo.RedisRepository
	emitter *emmit.Emitter
}

func NewLongPollHandler(repo *repo.RedisRepository, emmiter *emmit.Emitter) *longPollHandler {
	return &longPollHandler{emitter: emmiter, repo: repo}
}

func (l *longPollHandler) GetMessages(c echo.Context) error {
	currentUser := c.FormValue("username")
	timestampStr := c.QueryParam("timestamp")
	// Берём все сообщения если список пуст
	if timestampStr == "" || timestampStr == "null" {
		fmt.Println("timestampStr:", timestampStr)
		messages, err := l.repo.GetMessages(c.Request().Context())
		if err == nil && len(messages) > 0 {
			for _, msg := range messages {
				msg.IsSelf = msg.Username == currentUser
				if err := c.Render(http.StatusOK, "user-message", msg); err != nil {
					return err
				}
			}
			return nil
		}
	}

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
		Username:  username,
		Content:   message,
		Timestamp: time.Now().Unix(),
	}

	// Save to Redis
	if err := l.repo.SaveMessage(c.Request().Context(), msg); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not save message"})
	}

	l.emitter.Emit("newMessage", msg)

	return c.JSON(http.StatusOK, map[string]string{"message": message})
}
