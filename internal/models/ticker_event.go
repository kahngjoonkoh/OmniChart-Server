package models

import "github.com/google/uuid"

type TickerEvent struct {
	ID      uuid.UUID    `json:"id"`
	Ticker  Ticker `json:"ticker"`
	EventId uuid.UUID    `json:"event_id"`
	// Event *Event    `json:"event, omitempty"` // optional: populated if joined
}
