package models

import (
	"encoding/json"
	"time"
)

type ActivityLog struct {
	ID        int64           `json:"id"`
	ProjectID int64           `json:"project_id"`
	UserID    int64           `json:"user_id"`
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}
