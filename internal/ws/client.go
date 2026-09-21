package ws

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/anuzx/live-attendance-backend/internal/http/attendance"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
	sendBufferSize = 256
)

// Client = one connection. This is ws.user = {  userId: decoded.userId,  role: decoded.role}; .
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte   // outbound queue, drained by writePump
	done   chan struct{} // closed when this client is finished
	once   sync.Once
	userID uuid.UUID
	role   string
}

func newClient(h *Hub, conn *websocket.Conn, userID uuid.UUID, role string) *Client {
	return &Client{
		hub:    h,
		conn:   conn,
		send:   make(chan []byte, sendBufferSize),
		done:   make(chan struct{}),
		userID: userID,
		role:   role,
	}
}

func (c *Client) close() {
	c.once.Do(func() { close(c.done) }) // safe to call many times
}

// enqueue never blocks. Returns false if the client is gone or its buffer is full.
func (c *Client) enqueue(msg []byte) bool {
	select {
	case <-c.done:
		return false
	case c.send <- msg:
		return true
	default:
		return false
	}
}

func (c *Client) sendError(message string) {
	c.enqueue(encode("ERROR", ErrorData{Message: message}))
}

// readPump: one goroutine per connection, reads until the connection dies.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, raw, err := c.conn.ReadMessage() // blocks this goroutine only
		if err != nil {
			return // client left, timed out, or sent garbage
		}
		c.handleMessage(raw)
	}
}

// writePump: the ONLY place that writes to the socket.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.close()
		c.conn.Close()
	}()

	for {
		select {
		case msg := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}

		case <-ticker.C: // heartbeat so dead connections get noticed
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-c.done:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			c.conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			return
		}
	}
}

func (c *Client) handleMessage(raw []byte) {
	var msg IncomingMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		c.sendError("Invalid message format")
		return
	}

	switch msg.Event {
	case "ATTENDANCE_MARKED":
		c.handleAttendanceMarked(msg.Data)
	case "TODAY_SUMMARY":
		c.handleTodaySummary()
	case "MY_ATTENDANCE":
		c.handleMyAttendance()
	case "DONE":
		c.handleDone()
	default:
		c.sendError("Unknown event")
	}
}

func (c *Client) handleAttendanceMarked(raw json.RawMessage) {
	// 1. teacher only
	if c.role != "teacher" {
		c.sendError("Forbidden, teacher event only")
		return
	}

	var data AttendanceMarkedData
	if err := json.Unmarshal(raw, &data); err != nil ||
		data.StudentID == "" ||
		(data.Status != "present" && data.Status != "absent") {
		c.sendError("Invalid message format")
		return
	}

	studentID, err := uuid.Parse(data.StudentID)
	if err != nil {
		c.sendError("Invalid message format")
		return
	}

	// 2. needs an active session
	session, ok := c.hub.store.Get()
	if !ok {
		c.sendError("No active attendance session")
		return
	}

	// 3. the student must actually be enrolled in the active class
	enrolled, err := c.hub.attendanceSvc.IsStudentEnrolled(
		context.Background(), session.ClassID, studentID)
	if err != nil {
		log.Printf("check enrollment: %v", err)
		c.sendError("Internal server error")
		return
	}
	if !enrolled {
		c.sendError("Student is not enrolled in the active class")
		return
	}

	// 4. record it in memory
	if err := c.hub.store.MarkAttendance(data.StudentID, data.Status); err != nil {
		if errors.Is(err, attendance.ErrNoActiveSession) {
			c.sendError("No active attendance session")
			return
		}
		log.Printf("mark attendance: %v", err)
		c.sendError("Internal server error")
		return
	}

	// 5. broadcast to everyone
	c.hub.Broadcast(encode("ATTENDANCE_MARKED", data), c)
}

func (c *Client) handleTodaySummary() {
	// 1. teacher only
	if c.role != "teacher" {
		c.sendError("Forbidden, teacher event only")
		return
	}

	// 2. calculate from the in-memory session
	present, absent, total, err := c.hub.store.Summary()
	if err != nil {
		if errors.Is(err, attendance.ErrNoActiveSession) {
			c.sendError("No active attendance session")
			return
		}
		log.Printf("today summary: %v", err)
		c.sendError("Internal server error")
		return
	}

	// 3. broadcast to everyone
	c.hub.Broadcast(encode("TODAY_SUMMARY", SummaryData{
		Present: present,
		Absent:  absent,
		Total:   total,
	}), c)
}

func (c *Client) handleMyAttendance() {
	// 1. student only
	if c.role != "student" {
		c.sendError("Forbidden, student event only")
		return
	}

	// 2. look up this student's status
	status, err := c.hub.store.GetStatus(c.userID.String())
	if err != nil {
		if errors.Is(err, attendance.ErrNoActiveSession) {
			c.sendError("No active attendance session")
			return
		}
		log.Printf("my attendance: %v", err)
		c.sendError("Internal server error")
		return
	}

	// 3. unicast: only this client's own queue
	c.enqueue(encode("MY_ATTENDANCE", MyAttendanceData{Status: status}))
}

func (c *Client) handleDone() {
	// 1. teacher only
	if c.role != "teacher" {
		c.sendError("Forbidden, teacher event only")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	summary, err := c.hub.attendanceSvc.FinishSession(ctx)
	if err != nil {
		if errors.Is(err, attendance.ErrNoActiveSession) {
			c.sendError("No active attendance session")
			return
		}
		log.Printf("finish session: %v", err)
		c.sendError("Internal server error")
		return
	}

	// 7. broadcast to everyone
	c.hub.Broadcast(encode("DONE", DoneData{
		Message: "Attendance persisted",
		Present: summary.Present,
		Absent:  summary.Absent,
		Total:   summary.Total,
	}), c)
}
