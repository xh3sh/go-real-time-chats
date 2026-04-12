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
	"github.com/xh3sh/go-real-time-chats/internal/emmit"
)

type mockRenderer struct{}

func (m *mockRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return nil
}

func TestLongPollGetMessages(t *testing.T) {
	e := echo.New()
	e.Renderer = &mockRenderer{}

	t.Run("wait timeout", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?timestamp=123", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		repo := &mockRepo{}
		emitter := &mockEmitter{waitResOk: false}
		h := NewLongPollHandler(repo, emitter)

		err := h.GetMessages(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("receive message via wait", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?timestamp=123", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		repo := &mockRepo{}
		emitter := &mockEmitter{
			waitResOk:  true,
			waitResMsg: emmit.Message{Username: "user1", Content: "hi"},
		}
		h := NewLongPollHandler(repo, emitter)

		err := h.GetMessages(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestLongPollPostMessage(t *testing.T) {
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
		h := NewLongPollHandler(repo, emitter)

		if assert.NoError(t, h.PostMessage(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "newMessage", emitter.emittedEvent)
		}
	})
}
