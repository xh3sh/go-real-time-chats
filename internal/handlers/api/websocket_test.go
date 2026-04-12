package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/xh3sh/go-real-time-chats/internal/emmit"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(filename), "..", "..", "..")
	_ = os.Chdir(dir)
}

func TestWebSocketHandler(t *testing.T) {
	e := echo.New()
	repo := &mockRepo{}
	// Важно: создаем эмиттер с работающей подпиской (уже реализовано в mocks_test.go)
	emitter := &mockEmitter{}
	h := NewWebsocketHandler(repo, emitter)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := e.NewContext(r, w)
		_ = h.WebSocketHandler(c)
	}))
	defer srv.Close()

	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "?username=testuser"
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
	defer ws.Close()

	t.Run("client to server and echo back", func(t *testing.T) {
		testMsg := emmit.Message{Content: "ping"}
		err := ws.WriteJSON(testMsg)
		assert.NoError(t, err)

		// Ждем эхо-сообщения (сервер делает Emit при получении, и сам же его ловит через On)
		_, message, err := ws.ReadMessage()
		assert.NoError(t, err)
		assert.Contains(t, string(message), "ping")
		assert.Contains(t, string(message), "testuser")
	})

	t.Run("server to client broadcast", func(t *testing.T) {
		serverMsg := emmit.Message{Username: "system", Content: "broadcast"}
		emitter.Emit("newMessage", serverMsg)

		_, message, err := ws.ReadMessage()
		assert.NoError(t, err)
		assert.Contains(t, string(message), "broadcast")
		assert.Contains(t, string(message), "system")
	})
}
