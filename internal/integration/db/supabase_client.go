package db

import (
	"log"
	"os"
	"time"

	"github.com/supabase-community/supabase-go"
    "github.com/supabase-community/postgrest-go"

    "omnichart-server/internal/models"
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
        Ascending: false,  // false for descending order
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
