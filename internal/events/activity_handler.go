package events

import (
	"context"
	"encoding/json"
	"fluxera/internal/models"
)

type ActivityCreator interface {
	Create(ctx context.Context, projectID, userID int64, eventType string, payload json.RawMessage) (*models.ActivityLog, error)
}

type ActivityHandler struct {
	activity ActivityCreator
}

func NewActivityHandler(activity ActivityCreator) *ActivityHandler {
	return &ActivityHandler{activity: activity}
}

func (a *ActivityHandler) Handle(ctx context.Context, event models.Event) error {
	_, err := a.activity.Create(ctx, event.ProjectID, event.UserID, event.Type, event.Payload)
	return err
}
