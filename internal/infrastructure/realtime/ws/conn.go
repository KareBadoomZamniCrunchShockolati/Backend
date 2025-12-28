package ws

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Conn struct {
	ws   *websocket.Conn
	mu   sync.Mutex
	dead chan struct{}
}

func NewConn(ws *websocket.Conn) *Conn {
	_ = ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	ws.SetPongHandler(func(string) error {
		_ = ws.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	return &Conn{ws: ws, dead: make(chan struct{})}
}

func (c *Conn) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	select {
	case <-c.dead:
		return
	default:
		close(c.dead)
		_ = c.ws.Close()
	}
}

func (c *Conn) Send(payload []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	select {
	case <-c.dead:
		return
	default:
		_ = c.ws.WriteMessage(websocket.TextMessage, payload)
	}
}

func (c *Conn) ReadLoop(onClose func()) {
	defer func() {
		onClose()
		c.Close()
	}()
	for {
		if _, _, err := c.ws.ReadMessage(); err != nil {
			return
		}
	}
}
