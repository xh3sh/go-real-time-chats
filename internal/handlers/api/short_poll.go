package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xh3sh/go-real-time-chats/internal/emmit"
	"github.com/xh3sh/go-real-time-chats/internal/repo"
)

type ShortPoolHandler struct {
	repo    *repo.RedisRepository
	emitter *emmit.Emitter
}

func NewShortPollHandler(repo *repo.RedisRepository, emitter *emmit.Emitter) *ShortPoolHandler {
	return &ShortPoolHandler{
		repo:    repo,
		emitter: emitter,
	}
}

func (s *ShortPoolHandler) GetShortPollMessages(c echo.Context) error {
	currentUser := c.FormValue("username")
	if currentUser == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "currentUser is required"})
	}

	timestampStr := c.QueryParam("timestamp")
	var messages []emmit.Message
	var err error

	if timestampStr != "" && timestampStr != "null" {
		timestamp, parseErr := strconv.ParseInt(timestampStr, 10, 64)
		if parseErr != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid timestamp format"})
		}
		messages, err = s.repo.GetMessagesAfter(c.Request().Context(), timestamp)
	} else {
		messages, err = s.repo.GetMessages(c.Request().Context())
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not retrieve messages"})
	}

	for i := range messages {
		messages[i].IsSelf = (messages[i].Username == currentUser)
	}

	if len(messages) > 0 {
		for _, msg := range messages {
			if err := c.Render(http.StatusOK, "user-message", msg); err != nil {
				return err
			}
		}
		return nil
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *ShortPoolHandler) PostShortPollMessage(c echo.Context) error {
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

	if err := s.repo.SaveMessage(c.Request().Context(), msg); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not save message"})
	}

	s.emitter.Emit("newMessage", msg)

	return c.JSON(http.StatusOK, map[string]string{"message": message})
}



