package events

import (
	"context"
	"fluxera/internal/models"
	"fluxera/internal/service"
)

type ActivityHandler struct {
	activity *service.ActivityLogService
}

func NewActivityHandler(activity *service.ActivityLogService) *ActivityHandler {
	return &ActivityHandler{activity: activity}
}

func (a *ActivityHandler) Handle(ctx context.Context, event models.Event) error {
	_, err := a.activity.Create(ctx, event.ProjectID, event.UserID, event.Type, event.Payload)
	return err
}
