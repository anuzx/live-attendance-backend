package ws

//this file upgrades the http request to a websocket , verifies the token , creates a client ,registers it with the hub and starts the client's two go routines

import (
	"net/http"
	"time"

	"github.com/anuzx/live-attendance-backend/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wss = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Browsers send an Origin header and gorilla rejects cross-origin by default.
	// Fine for development; restrict this in production.
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Hub) ServeWS(c *gin.Context) {
	userID, role, err := auth.ParseToken(c.Query("token"))

	if err != nil {
		// invalid token: upgrade only to send ERROR and close, as before
		conn, upErr := wss.Upgrade(c.Writer, c.Request, nil)
		if upErr != nil {
			return
		}
		conn.WriteMessage(websocket.TextMessage,
			encode("ERROR", ErrorData{Message: "Unauthorized or invalid token"}))
		conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "unauthorized"),
			time.Now().Add(time.Second),
		)
		conn.Close()
		return
	}

	// register first, then upgrade: the 101 response (and therefore the
	// client's "open" event) is only written after registration, so a
	// connection is always broadcast-ready by the time it can send messages.
	// Queued broadcasts are buffered in c.send until writePump starts.
	client := newClient(h, nil, userID, role)
	h.register <- client

	conn, err := wss.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.unregister <- client
		return
	}

	client.conn = conn
	go client.writePump()
	go client.readPump()
}
