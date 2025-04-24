package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xh3sh/go-real-time-chats/internal/emmit"
)

var shortPollMessages []MessageWithTimestamp
var lastRequestTimestamps = make(map[string]time.Time)

type MessageWithTimestamp struct {
	emmit.Message
	Timestamp time.Time
}

func GetShortPollMessages(c echo.Context) error {
	currentUser := c.FormValue("username")
	if currentUser == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "currentUser is required"})
	}

	// Get the last request timestamp for the user
	lastTimestamp, exists := lastRequestTimestamps[currentUser]
	if !exists {
		lastTimestamp = time.Time{}
	}

	lastRequestTimestamps[currentUser] = time.Now()

	cleanupInactiveUsers()

	cleanupOldMessages()

	var filteredMessages []emmit.Message
	for _, msg := range shortPollMessages {
		if msg.Timestamp.After(lastTimestamp) {
			msg.IsSelf = (msg.Username == currentUser)
			filteredMessages = append(filteredMessages, msg.Message)
		}
	}

	if len(filteredMessages) > 0 {
		for _, msg := range filteredMessages {
			if err := c.Render(http.StatusOK, "user-message", msg); err != nil {
				return err
			}
		}
		return nil
	}

	return c.NoContent(http.StatusNoContent)
}

func PostShortPollMessage(c echo.Context) error {
	message := c.FormValue("message")
	username := c.FormValue("username")

	if username == "" || message == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username and message cannot be empty"})
	}

	msg := MessageWithTimestamp{
		Message: emmit.Message{
			Username: username,
			Content:  message,
		},
		Timestamp: time.Now(),
	}

	shortPollMessages = append(shortPollMessages, msg)

	cleanupInactiveUsers()

	return c.JSON(http.StatusOK, map[string]string{"message": message})
}

func cleanupOldMessages() {
	threshold := time.Now().Add(-1 * time.Minute)
	var updatedMessages []MessageWithTimestamp
	for _, msg := range shortPollMessages {
		if msg.Timestamp.After(threshold) {
			updatedMessages = append(updatedMessages, msg)
		}
	}
	shortPollMessages = updatedMessages
}

func cleanupInactiveUsers() {
	threshold := time.Now().Add(-5 * time.Minute)
	for user, lastRequest := range lastRequestTimestamps {
		if lastRequest.Before(threshold) {
			delete(lastRequestTimestamps, user)
		}
	}
}
