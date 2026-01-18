package emmit

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Username  string
	Content   string
	Timestamp int64
	IsSelf    bool
}

type Listener func(message Message)

type Emitter struct {
	mu        sync.Mutex
	listeners map[string]map[string]Listener
}

func New() *Emitter {
	return &Emitter{
		listeners: make(map[string]map[string]Listener),
	}
}

func (e *Emitter) ListenerCount() int {
	l, ok := e.listeners["newMessage"]
	if ok {
		return len(l)
	}
	return 0
}

// AddListener добавляет слушателя для события
func (e *Emitter) AddListener(event string, listener Listener) (key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.listeners[event] == nil {
		e.listeners[event] = make(map[string]Listener)
	}
	key = uuid.NewString()
	e.listeners[event][key] = listener
	return key
}

// RemoveListener удаляет слушателя по ключу
func (e *Emitter) RemoveListener(event, key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if listeners, exists := e.listeners[event]; exists {
		delete(listeners, key)
		if len(listeners) == 0 {
			delete(e.listeners, event)
		}
	}
}

// On создает постоянного слушателя для события
func (e *Emitter) On(event string, listener Listener) (off func()) {
	key := e.AddListener(event, listener)
	return func() {
		e.RemoveListener(event, key)
	}
}

// Emit отправляет сообщение всем слушателям события
func (e *Emitter) Emit(event string, message Message) {
	e.mu.Lock()
	listeners := e.listeners[event]
	e.mu.Unlock()

	for _, listener := range listeners {
		go listener(message)
	}
}

// Wait ожидает сообщение в канале с таймаутом
func (e *Emitter) Wait(event string, timeout time.Duration) (Message, bool) {
	ch := make(chan Message, 1)
	key := e.AddListener(event, func(msg Message) { ch <- msg })

	select {
	case msg := <-ch:
		e.RemoveListener(event, key)
		return msg, true
	case <-time.After(timeout):
		e.RemoveListener(event, key)
		return Message{}, false
	}
}
