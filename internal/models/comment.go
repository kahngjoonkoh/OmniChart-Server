package models

import (
	"github.com/google/uuid"
	"time"
)


type Comment struct {
	ID         		uuid.UUID `json:"id"`
	UserID    		string    `json:"user_id"`
	TickerEventID   string    `json:"ticker_event_id"` // Check ticker and event to get proper comment seciton
	Content   		string    `json:"content"`
	CreatedAt 		time.Time `json:"created_at"`
}
