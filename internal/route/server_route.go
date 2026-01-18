package route

import (
	"github.com/labstack/echo/v4"
	"github.com/xh3sh/go-real-time-chats/internal/emmit"
	"github.com/xh3sh/go-real-time-chats/internal/handlers/api"
	"github.com/xh3sh/go-real-time-chats/internal/repo"
)

func serverRoute(c *echo.Group, repoRDB *repo.RedisRepository) {
	emitter := emmit.New()

	longPollHandler := api.NewLongPollHandler(repoRDB, emitter)
	c.GET("/long-poll-messages", longPollHandler.GetMessages)
	c.POST("/long-poll-messages", longPollHandler.PostMessage)

	shortPoolHandler := api.NewShortPollHandler(repoRDB, emitter)
	c.GET("/short-poll-messages", shortPoolHandler.GetShortPollMessages)
	c.POST("/short-poll-messages", shortPoolHandler.PostShortPollMessage)

	wsHandler := api.NewWebsocketHandler(repoRDB, emitter)
	c.GET("/websocket", wsHandler.WebSocketHandler)

	sseHandler := api.NewSSEHandler(repoRDB, emitter)
	c.GET("/sse", sseHandler.SSEHandler)
	c.POST("/sse-message", sseHandler.SSEMessageHandler)
	c.GET("/users-online", sseHandler.UserOnline)

	api.UsersOnlineTicker(emitter)
}
