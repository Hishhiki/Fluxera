package models

import (
	"encoding/json"
	"time"
)

const (
	EventProjectCreated    = "project.created"
	EventTaskCreated       = "task.created"
	EventTaskUpdated       = "task.updated"
	EventTaskStatusChanged = "task.status_changed"
	EventCommentCreated    = "comment.created"
)

type Event struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	ProjectID int64           `json:"project_id"`
	UserID    int64           `json:"user_id"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}
