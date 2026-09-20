package ws

import "encoding/json"

//	{
//	  "event": "EVENT_NAME",
//	  "data": { ... }
//	}
type IncomingMessage struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

type OutgoingMessage struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

type AttendanceMarkedData struct {
	StudentID string `json:"studentId"`
	Status    string `json:"status"`
}
type ErrorData struct {
	Message string `json:"message"`
}

func encode(event string, data any) []byte {
	b, _ := json.Marshal(OutgoingMessage{Event: event, Data: data}) // these types always marshal
	return b
}
