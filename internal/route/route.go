package route

import (
	"github.com/labstack/echo/v4"
	"github.com/xh3sh/go-real-time-chats/internal/templates"
)

type RouteHandler struct {
	echo *echo.Echo
}

func New(e *echo.Echo) *RouteHandler {
	return &RouteHandler{
		echo: e,
	}
}

func (r *RouteHandler) InitRoute(tmpl *templates.Templates) error {
	r.echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("templates", tmpl)
			return next(c)
		}
	})

	serverRoute(r.echo.Group("/api"))

	webRoute(r.echo.Group(""))

	return nil
}
