package models

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID          uuid.UUID `json:"id"`
	Timestamp   time.Time `json:"timestamp"` // Use timestamptz in Postgres
	Title       string    `json:"title"`
	SourceUrl   string    `json:"source_url"`
	Content     string    `json:"content"`
	EventTypeID int       `json:"event_type_id"`
	// EventType   *EventType `json:"event_type,omitempty"` // optional: populated if joined
}
