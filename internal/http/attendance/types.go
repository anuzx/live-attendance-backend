package attendance

type StartSessionRequest struct {
	ClassID string `json:"classId" binding:"required"`
}

type StartSessionResponse struct {
	ClassID   string `json:"classId"`
	StartedAt string `json:"startedAt"`
}