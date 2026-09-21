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
type SummaryData struct {
	Present int `json:"present"`
	Absent  int `json:"absent"`
	Total   int `json:"total"`
}
type MyAttendanceData struct {
	Status string `json:"status"`
}

type DoneData struct {
	Message string `json:"message"`
	Present int    `json:"present"`
	Absent  int    `json:"absent"`
	Total   int    `json:"total"`
}

func encode(event string, data any) []byte {
	b, _ := json.Marshal(OutgoingMessage{Event: event, Data: data})
	return b
}
