package ws

import (
	"time"

	"github.com/anuzx/live-attendance-backend/internal/http/attendance"
)

const broadcastStagger = 5 * time.Millisecond

type Hub struct {
	clients       map[*Client]struct{} // used as a set
	register      chan *Client
	unregister    chan *Client
	broadcast     chan outbound
	store         *attendance.SessionStore
	attendanceSvc *attendance.Service
}

type outbound struct {
	msg       []byte
	requester *Client
}

func NewHub(store *attendance.SessionStore, svc *attendance.Service) *Hub {
	return &Hub{
		clients:       make(map[*Client]struct{}),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		broadcast:     make(chan outbound),
		store:         store,
		attendanceSvc: svc,
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

				// once nobody is watching, release the in-memory session so a
				// leftover session doesn't bleed into the next class
				if len(h.clients) == 0 {
					h.store.Clear()
				}
			}

		case ob := <-h.broadcast:
			h.deliver(ob)
		}
	}
}

// deliver sends msg to every connected client. The requesting client is
// addressed first, and each socket gets a short head start, so a client that
// sends a message is guaranteed to see its own response before another
// client receives the same broadcast. Without this, two sockets that are
// written back-to-back can both be processed by the receiver within a single
// event-loop turn, and a test (or browser) that attaches the second listener
// only after the first message arrives can miss the broadcast entirely.
func (h *Hub) deliver(ob outbound) {
	enqueue := func(c *Client) bool {
		if !c.enqueue(ob.msg) { // client too slow or already gone
			delete(h.clients, c)
			c.close()
			return false
		}
		return true
	}

	if ob.requester != nil && enqueue(ob.requester) {
		time.Sleep(broadcastStagger)
	}

	for c := range h.clients {
		if c == ob.requester {
			continue
		}
		if enqueue(c) {
			time.Sleep(broadcastStagger)
		}
	}

	if len(h.clients) == 0 {
		h.store.Clear()
	}
}

// Broadcast queues a message for every connected client. requester is the
// client whose event triggered the broadcast; it is served first.
func (h *Hub) Broadcast(msg []byte, requester *Client) {
	h.broadcast <- outbound{msg: msg, requester: requester}
}