package events

import (
	"context"
	"fluxera/internal/models"
	"log"
	"time"

	"github.com/google/uuid"
)

type Dispatcher struct {
	events chan models.Event
}
type Handler interface {
	Handle(ctx context.Context, event models.Event) error
}

func NewDispatcher(buffer int) *Dispatcher {
	return &Dispatcher{
		events: make(chan models.Event, buffer),
	}
}

func (d *Dispatcher) Publish(ctx context.Context, event models.Event) error {
	if event.ID == "" {
		event.ID = uuid.NewString()
	}

	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}
	select {
	case d.events <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (d *Dispatcher) Start(ctx context.Context, handler Handler) {
	go func() {
		for {
			select {
			case event, ok := <-d.events:
				if !ok {
					return
				}
				if err := handler.Handle(ctx, event); err != nil {
					log.Printf("failed to handle event: %v", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
