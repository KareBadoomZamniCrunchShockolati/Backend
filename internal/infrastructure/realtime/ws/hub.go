package ws

import "sync"

type Hub struct {
	mu    sync.RWMutex
	conns map[uint]map[*Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{conns: make(map[uint]map[*Conn]struct{})}
}

func (h *Hub) Add(userID uint, c *Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.conns[userID] == nil {
		h.conns[userID] = make(map[*Conn]struct{})
	}
	h.conns[userID][c] = struct{}{}
}

func (h *Hub) Remove(userID uint, c *Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if m := h.conns[userID]; m != nil {
		delete(m, c)
		if len(m) == 0 {
			delete(h.conns, userID)
		}
	}
}

func (h *Hub) Send(userID uint, payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.conns[userID] {
		c.Send(payload)
	}
}
