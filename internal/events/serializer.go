package events

import (
	"encoding/json"
	"fluxera/internal/models"
)

func SerializeEvent(event models.Event) ([]byte, error) {
	return json.Marshal(event)
}
