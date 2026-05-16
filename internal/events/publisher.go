package events

import (
	"context"
	"fluxera/internal/models"
)

type Publisher interface {
	Publish(ctx context.Context, event models.Event) error
}
