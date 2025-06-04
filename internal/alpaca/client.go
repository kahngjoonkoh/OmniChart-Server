package alpacaApi

import (
	"log"
	"os"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

var Client *alpaca.Client

func Init() {
	apiKey := os.Getenv("APCA_API_KEY_ID")
	apiSecret := os.Getenv("APCA_API_SECRET_KEY")
	baseURL := os.Getenv("APCA_API_BASE_URL") // Optional, defaults to paper trading

	if apiKey == "" || apiSecret == "" {
		log.Fatal("APCA_API_KEY_ID and APCA_API_SECRET_KEY must be set")
	}

	if baseURL == "" {
		baseURL = "https://paper-api.alpaca.markets"
	}

	client := alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   baseURL,
	})

	// Check if the credentials are valid
	if _, err := client.GetAccount(); err != nil {
		log.Fatalf("Failed to authenticate with Alpaca API: %v", err)
	}

	Client = client
	log.Println("Successfully initialized Alpaca client.")
}
