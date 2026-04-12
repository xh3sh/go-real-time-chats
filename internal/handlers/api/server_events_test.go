package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSSEMessageHandler(t *testing.T) {
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		f := make(url.Values)
		f.Set("username", "user1")
		f.Set("content", "hello")

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		repo := &mockRepo{}
		emitter := &mockEmitter{}
		h := NewSSEHandler(repo, emitter)

		if assert.NoError(t, h.SSEMessageHandler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "newMessage", emitter.emittedEvent)
		}
	})

	t.Run("empty message", func(t *testing.T) {
		f := make(url.Values)
		f.Set("username", "user1")
		f.Set("content", "")

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := NewSSEHandler(&mockRepo{}, &mockEmitter{})
		err := h.SSEMessageHandler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
