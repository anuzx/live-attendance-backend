package ws

import "github.com/anuzx/live-attendance-backend/internal/http/attendance"

type Hub struct {
	clients    map[*Client]struct{} // used as a set
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	store      *attendance.SessionStore
}

func NewHub(store *attendance.SessionStore) *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		store:      store,
	}
}

// Run must be started once with `go hub.Run()`. It is the ONLY goroutine that touches h.clients.
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = struct{}{}

		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				c.close()
			}

		case msg := <-h.broadcast:
			for c := range h.clients {
				if !c.enqueue(msg) { // client too slow or already gone
					delete(h.clients, c)
					c.close()
				}
			}
		}
	}
}

func (h *Hub) Broadcast(msg []byte) {
	h.broadcast <- msg
}
