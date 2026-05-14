package models

import "encoding/json"

type Event struct {
	ProjectID int64           `json:"project_id"`
	UserID    int64           `json:"user_id"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
}
