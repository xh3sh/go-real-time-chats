package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/xh3sh/go-real-time-chats/internal/route"
	"github.com/xh3sh/go-real-time-chats/internal/templates"
)

const SERVER_PORT = ":80"

func main() {
	e := echo.New()

	tmpl := templates.NewTemplates()
	e.Renderer = tmpl

	e.Use(middleware.Logger())

	routes := route.New(e)
	routes.InitRoute(tmpl)

	e.Logger.Fatal(e.Start(SERVER_PORT))
}
