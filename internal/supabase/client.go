package supabase

import (
	"log"
	"os"
	"time"

	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"

	"omnichart-server/internal/models"

	"fmt"
	"encoding/json"
	"github.com/google/uuid"
)

var Client *supabase.Client

func Init() {
    var err error

	DB_URL := os.Getenv("SUPABASE_URL")
	DB_KEY := os.Getenv("SUPABASE_ANON_KEY")

	if DB_URL == "" || DB_KEY == "" {
		log.Fatal("SUPABASE_URL and SUPABASE_ANON_KEY must be set")
	}

	Client, err = supabase.NewClient(DB_URL, DB_KEY, nil)
	if err != nil {
		log.Fatalf("Failed to create Supabase client: %v", err)
	}
	log.Println("Successfully Initialized Supabase client.")
}

func GetEvents(ticker string, from time.Time, to time.Time, limit int) ([]models.Event, error) {
	var events []models.Event

	fromStr := from.Format(time.RFC3339)
	toStr := to.Format(time.RFC3339)

	orderOpts := &postgrest.OrderOpts{
		Ascending: false, // false for descending order
	}

	data, count, err := Client.From("ticker_event").
		Select("event!inner(id,timestamp,title,source_url,content,event_type)", "exact", false).
		Eq("ticker.ticker", ticker).
		Gte("event.timestamp", fromStr).
		Lte("event.timestamp", toStr).
		Order("event.timestamp", orderOpts).
		Limit(limit, "").Execute()

	if err != nil {
		return nil, err
	}

	log.Println(data, count)

	// // extract events from results.
	// events = make([]models.Event, count)
	// for i, r := range data {
	//     events[i] = r.Event
	// }

	return events, nil
}

func GetSupabaseClient() *supabase.Client {
	return Client
}

func AddComment(tickerEventID, userID, content string) (*models.Comment, error) {
	client := GetSupabaseClient()

	insert := map[string]interface{}{
		"ticker_event_id": tickerEventID,
		"user_id":         userID,
		"content":         content,
	}

	resp, count, err := client.From("comments").Insert(insert, false, "", "representation", "").Execute()

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, fmt.Errorf("no rows inserted")
	}

	// Unmarshal resp (JSON) into []models.Comment
	var inserted []models.Comment
	err = json.Unmarshal(resp, &inserted)
	if err != nil {
		return nil, err
	}

	return &inserted[0], nil
}

func GetComments(tickerEventID string) ([]models.Comment, error) {
	client := GetSupabaseClient()

	// Compose order options for ascending order
	orderOpts := &postgrest.OrderOpts{Ascending: true}

	// Select params: columns, head ("" for no head), count (false)
	resp, count, err := client.From("comments_with_usernames").
		Select("id, user_id, username, ticker_event_id, content, created_at", "", false).
		Eq("ticker_event_id", tickerEventID).
		Order("created_at", orderOpts).
		Execute()

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return []models.Comment{}, nil
	}

	var comments []models.Comment
	err = json.Unmarshal(resp, &comments)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func AddTickerEvent(ticker string, eventID string, startIndex, endIndex int) (*models.TickerEvent, error) {
	client := GetSupabaseClient()

	// Convert eventID to uuid.UUID to validate format
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return nil, fmt.Errorf("invalid event ID format: %v", err)
	}

	insert := map[string]interface{}{
		"ticker":      ticker,
		"event_id":    eventUUID,
		"start_index": startIndex,
		"end_index":   endIndex,
	}

	resp, count, err := client.From("ticker_events").Insert(insert, false, "", "representation", "").Execute()
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, fmt.Errorf("no rows inserted")
	}

	var inserted []models.TickerEvent
	err = json.Unmarshal(resp, &inserted)
	if err != nil {
		return nil, err
	}

	return &inserted[0], nil
}

func GetTickerEvents(ticker string) ([]models.TickerEvent, error) {
	client := GetSupabaseClient()

	resp, count, err := client.From("ticker_events").
		Select("id, ticker, event_id, start_index, end_index", "", false).
		Eq("ticker", ticker).
		Execute()

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return []models.TickerEvent{}, nil
	}

	var tickerEvents []models.TickerEvent
	err = json.Unmarshal(resp, &tickerEvents)
	if err != nil {
		return nil, err
	}

	return tickerEvents, nil
}
