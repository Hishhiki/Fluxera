package models

import (
	"encoding/json"
	"time"
)

const (
	EventTaskCreated       = "task.created"
	EventTaskStatusChanged = "task.status_changed"
	EventCommentCreated    = "comment.created"
	EventTaskUpdated       = "task.updated"
	EventProjectCreated    = "project.created"
)

type Event struct {
	ID        string          `json:"id"`
	ProjectID int64           `json:"project_id"`
	UserID    int64           `json:"user_id"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}
