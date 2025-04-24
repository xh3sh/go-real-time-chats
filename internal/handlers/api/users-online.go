package api

import (
	"log"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xh3sh/go-real-time-chats/internal/emmit"
)

const USERS_TICKER_TIME = 30

// UsersOnlineTicker отправляет количество слушаетелй в UserOnline хэндлер
func UsersOnlineTicker(emitter *emmit.Emitter) {
	ticker := time.NewTicker(USERS_TICKER_TIME * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				emitter.Emit("usersOnline", emmit.Message{Content: strconv.Itoa(emitter.ListenerCount())})
			}
		}
	}()
}

// UserOnline отображает количество слушателей онлайн
func (s *sseHandler) UserOnline(c echo.Context) error {
	client := &SSEClient{
		Channel: make(chan emmit.Message),
	}

	sseClientsMux.Lock()
	sseClients[client] = struct{}{}
	sseClientsMux.Unlock()

	defer func() {
		sseClientsMux.Lock()
		delete(sseClients, client)
		sseClientsMux.Unlock()
		close(client.Channel)
	}()

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")

	if _, err := c.Response().Write([]byte("data:" + strconv.Itoa(s.emitter.ListenerCount()) + "\n\n")); err != nil {
		log.Println(err)
	}
	c.Response().Flush()

	off := s.emitter.On("usersOnline", func(msg emmit.Message) {
		if _, err := c.Response().Write([]byte("data:" + msg.Content + "\n\n")); err != nil {
			log.Println(err)
		}
		c.Response().Flush()
	})

	defer off()

	<-c.Request().Context().Done()
	return nil
}
