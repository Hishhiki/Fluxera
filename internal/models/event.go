package models

import "encoding/json"

const (
	EventTaskCreated       = "task.created"
	EventTaskStatusChanged = "task.status_changed"
	EventCommentCreated    = "comment.created"
)

type Event struct {
	ProjectID int64           `json:"project_id"`
	UserID    int64           `json:"user_id"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
}
