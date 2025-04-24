package route

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xh3sh/go-real-time-chats/internal/handlers/web"
)

func webRoute(c *echo.Group) {
	fs := http.FileServer(http.Dir("web/static"))
	c.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", fs)))

	c.GET("/", web.HandleHome)
	c.GET("/ws", web.HandleWS)
	c.GET("/event", web.HandleEvent)
	c.GET("/load-template", web.HandleLoadTemplate)
	c.GET("/long-poll-login", web.LongPollLogin)
	c.GET("/ws-login", web.WSLogin)
	c.GET("/sse-login", web.SSELogin)
	c.GET("/short-poll-login", web.ShortPollLogin)
}
