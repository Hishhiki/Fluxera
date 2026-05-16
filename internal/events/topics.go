package events

import (
	"errors"
	"fluxera/internal/models"
)

var ErrUnknownEventType = errors.New("unknown event type")

func TopicForEventType(eventType string) (string, error) {
	switch eventType {
	case models.EventProjectCreated:
		return "project.created", nil
	case models.EventTaskCreated:
		return "task.created", nil
	case models.EventTaskUpdated:
		return "task.updated", nil
	case models.EventTaskStatusChanged:
		return "task.status_changed", nil
	case models.EventCommentCreated:
		return "comment.created", nil
	default:
		return "", ErrUnknownEventType
	}
}

func TopicForEvent(event models.Event) (string, error) {
	return TopicForEventType(event.Type)
}
