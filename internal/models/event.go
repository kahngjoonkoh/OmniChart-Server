package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID  `json:"id"`
	Timestamp   time.Time  `json:"timestamp"` // Use timestamptz in Postgres
	Title       string     `json:"title"`
	SourceUrl   string     `json:"source_url"`
	Content     string     `json:"content"`
	EventTypesID int        `json:"event_types_id"`
	EventType   *EventType `json:"event_type,omitempty"`
}
