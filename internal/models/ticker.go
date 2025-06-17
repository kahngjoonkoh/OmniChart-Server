package models

import "time"

type Ticker struct {
	Ticker string `json:"ticker"`
	Name   string `json:"name"`
	LastUpdated time.Time `json:last_updated`
}
