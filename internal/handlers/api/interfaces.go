package api

import (
	"context"
	"time"

	"github.com/xh3sh/go-real-time-chats/internal/emmit"
)

type MessageRepository interface {
	SaveMessage(ctx context.Context, msg emmit.Message) error
	GetMessages(ctx context.Context) ([]emmit.Message, error)
	GetMessagesAfter(ctx context.Context, timestamp int64) ([]emmit.Message, error)
}

type EventEmitter interface {
	Emit(event string, message emmit.Message)
	Wait(event string, timeout time.Duration) (emmit.Message, bool)
	On(event string, listener emmit.Listener) (off func())
	ListenerCount() int
}
