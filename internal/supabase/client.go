package supabase

import (
	"log"
	"os"
	"sort"

	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"

	"omnichart-server/internal/models"

	"encoding/json"
	"fmt"
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

func GetEvents(ticker string) ([]models.TickerEvent, error) {
	var results []models.TickerEvent
	selectQuery := "id,ticker,event_id,start_index,end_index,events!inner(id,timestamp,title,source_url,content,event_types_id)"

	data, _, err := Client.From("ticker_event").
		Select(selectQuery, "exact", false).
		Eq("ticker", ticker).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("supabase query error for ticker '%s': %w", ticker, err)
	}

	err = json.Unmarshal(data, &results)

	log.Println(results, err)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling Supabase data: %w", err)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Event.Timestamp.After(results[j].Event.Timestamp) // Descending order
	})

	return results, nil
}

func GetSupabaseClient() *supabase.Client {
	return Client
}

func AddComment(tickerEventID string, userID string, content string) (*models.Comment, error) {
	client := GetSupabaseClient()

	insert := map[string]interface{}{
		"content":         content,
		"ticker_event_id": tickerEventID,
		"user_id":         userID,
	}

	resp, _, err := client.
		From("comments").
		Insert(insert, false, "representation", "", "").Execute()

	if err != nil {
		return nil, err
	}

	var inserted []models.Comment
	err = json.Unmarshal(resp, &inserted)
	if err != nil {
		return nil, err
	}

	if len(inserted) == 0 {
		return nil, fmt.Errorf("no rows inserted")
	}

	return &inserted[0], nil
}

func GetComments(tickerEventID string) ([]models.Comment, error) {
	client := GetSupabaseClient()

	// Compose order options for ascending order
	orderOpts := &postgrest.OrderOpts{Ascending: true}

	// Select params: columns, head ("" for no head), count (false)
	resp, count, err := client.From("comments").
		Select("*", "exact", false).
		Eq("ticker_event_id", tickerEventID).
		Order("timestamp", orderOpts).
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
