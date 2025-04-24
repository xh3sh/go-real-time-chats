package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleHome(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func HandleWS(c echo.Context) error {
	return c.Render(http.StatusOK, "ws", nil)
}

func HandleEvent(c echo.Context) error {
	return c.Render(http.StatusOK, "event", nil)
}

func HandleLoadTemplate(c echo.Context) error {

	templateName := c.QueryParam("value")

	switch templateName {
	case "":
		return c.Render(http.StatusOK, "info-box", nil)
	case "short-poll":
		return c.Render(http.StatusOK, "short-poll", nil)
	case "long-poll":
		return c.Render(http.StatusOK, "long-poll", nil)
	case "ws":
		return c.Render(http.StatusOK, "ws", nil)
	case "sse":
		return c.Render(http.StatusOK, "sse", nil)
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":   "Invalid template",
			"message": "The requested template does not exist. Please select a valid template.",
		})
	}
}

func LongPollLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "enable-polling", nil)
}

func ShortPollLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "enable-short-polling", nil)
}

func WSLogin(c echo.Context) error {
	username := c.QueryParam("username")
	wsData := struct {
		Username string
	}{
		Username: username,
	}
	return c.Render(http.StatusOK, "ws-enable", wsData)
}

func SSELogin(c echo.Context) error {
	username := c.QueryParam("username")
	sseData := struct {
		Username string
	}{
		Username: username,
	}
	return c.Render(http.StatusOK, "sse-enable", sseData)
}
