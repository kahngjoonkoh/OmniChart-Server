package models

import "github.com/google/uuid"

type TickerEvent struct {
	ID      uuid.UUID    `json:"id"`
	Ticker  string `json:"ticker"`
	EventId uuid.UUID    `json:"event_id"`
	StartIndex int       `json:"start_index"`
	EndIndex   int       `json:"end_index"`
	// Event *Event    `json:"event, omitempty"` // optional: populated if joined
}
