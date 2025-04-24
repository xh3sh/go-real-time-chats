package route

import (
	"github.com/labstack/echo/v4"
	"github.com/xh3sh/go-real-time-chats/internal/emmit"
	"github.com/xh3sh/go-real-time-chats/internal/handlers/api"
)

func serverRoute(c *echo.Group) {
	emitter := emmit.New()

	longPollHandler := api.NewLongPollHandler(emitter)
	c.GET("/long-poll-messages", longPollHandler.GetMessages)
	c.POST("/long-poll-messages", longPollHandler.PostMessage)

	c.GET("/short-poll-messages", api.GetShortPollMessages)
	c.POST("/short-poll-messages", api.PostShortPollMessage)

	wsHandler := api.NewWebsocketHandler(emitter)
	c.GET("/websocket", wsHandler.WebSocketHandler)

	sseHandler := api.NewSSEHandler(emitter)
	c.GET("/sse", sseHandler.SSEHandler)
	c.POST("/sse-message", sseHandler.SSEMessageHandler)
	c.GET("/users-online", sseHandler.UserOnline)

	api.UsersOnlineTicker(emitter)
}
