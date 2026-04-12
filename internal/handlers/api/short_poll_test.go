package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockRendererShort struct{}

func (m *mockRendererShort) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return nil
}

func TestGetShortPollMessages(t *testing.T) {
	e := echo.New()
	e.Renderer = &mockRendererShort{}

	t.Run("missing username", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?timestamp=123", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := NewShortPollHandler(&mockRepo{}, &mockEmitter{})
		err := h.GetShortPollMessages(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid timestamp", func(t *testing.T) {
		f := make(url.Values)
		f.Set("username", "testuser")

		req := httptest.NewRequest(http.MethodGet, "/?timestamp=invalid", strings.NewReader(f.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := NewShortPollHandler(&mockRepo{}, &mockEmitter{})
		err := h.GetShortPollMessages(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestPostShortPollMessage(t *testing.T) {
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		f := make(url.Values)
		f.Set("username", "user1")
		f.Set("message", "hello")

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		repo := &mockRepo{}
		emitter := &mockEmitter{}
		h := NewShortPollHandler(repo, emitter)

		if assert.NoError(t, h.PostShortPollMessage(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "newMessage", emitter.emittedEvent)
			assert.Equal(t, "user1", emitter.emittedMsg.Username)
			assert.Equal(t, "hello", emitter.emittedMsg.Content)
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		f := make(url.Values)
		f.Set("username", "")

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := NewShortPollHandler(&mockRepo{}, &mockEmitter{})
		err := h.PostShortPollMessage(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
