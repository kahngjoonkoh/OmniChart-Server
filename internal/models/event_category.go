package models

type EventCategory struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
	Name     string `json:"name"`
}
