package models

type EventType struct {
	ID         int    `json:"id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	CategoryId int    `json:"category_id"`
	// Category   *EventCategory `json:"event_category,omitempty"` // optional: populated if joined
}
