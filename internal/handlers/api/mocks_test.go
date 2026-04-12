package api

import (
	"context"
	"time"

	"github.com/xh3sh/go-real-time-chats/internal/emmit"
)

type mockRepo struct {
	saveMessageErr    error
	getMessagesRes    []emmit.Message
	getMessagesErr    error
	getMessagesAfter  []emmit.Message
	getMessagesAfterE error
}

func (m *mockRepo) SaveMessage(ctx context.Context, msg emmit.Message) error {
	return m.saveMessageErr
}

func (m *mockRepo) GetMessages(ctx context.Context) ([]emmit.Message, error) {
	return m.getMessagesRes, m.getMessagesErr
}

func (m *mockRepo) GetMessagesAfter(ctx context.Context, ts int64) ([]emmit.Message, error) {
	return m.getMessagesAfter, m.getMessagesAfterE
}

type mockEmitter struct {
	emittedEvent string
	emittedMsg   emmit.Message
	waitResMsg   emmit.Message
	waitResOk    bool
	listeners    map[string]emmit.Listener
}

func (m *mockEmitter) Emit(event string, msg emmit.Message) {
	m.emittedEvent = event
	m.emittedMsg = msg
	if l, ok := m.listeners[event]; ok {
		l(msg)
	}
}

func (m *mockEmitter) Wait(event string, timeout time.Duration) (emmit.Message, bool) {
	return m.waitResMsg, m.waitResOk
}

func (m *mockEmitter) On(event string, listener emmit.Listener) (off func()) {
	if m.listeners == nil {
		m.listeners = make(map[string]emmit.Listener)
	}
	m.listeners[event] = listener
	return func() { delete(m.listeners, event) }
}

func (m *mockEmitter) ListenerCount() int {
	return 0
}
