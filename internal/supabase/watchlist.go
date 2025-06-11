package supabase

import (
	"encoding/json"
	"errors"

	"omnichart-server/internal/models"
)

// Add a ticker to user's watchlist
// Require an authenticated user to perform this action
func AddTickerToWatchlist(ticker string) error {
	// Ensure the current user is authenticated
	resp, _ := Client.Auth.GetUser()
	if resp == nil {
		return errors.New("User not logged in")
	}

	// Insert new entry to watchlists table
	req := map[string]interface{}{
		"ticker": ticker,
	}
	_, _, err := Client.From("watchlists").
		Insert(req, false, "ignore", "", "").
		Execute()
	return err
}

// Retrieve the tickers in user's watchlist
// Require an authenticated user to perform this action
func GetWatchlist() ([]string, error) {
	// Ensure the current user is authenticated
	userResp, _ := Client.Auth.GetUser()
	if userResp == nil {
		return nil, errors.New("User not logged in")
	}

	// Look up tickers in user's watchlist
	resp, _, err := Client.From("watchlists").
		Select("ticker", "exact", false).
		Eq("user_id", userResp.User.ID.String()).
		Execute()
	if err != nil {
		return nil, err
	}

	// Decode the response data
	var tickersJson []models.Watchlist
	err = json.Unmarshal(resp, &tickersJson)
	if err != nil {
		return nil, err
	}
	tickers := make([]string, len(tickersJson))
	for i, ticker := range tickersJson {
		tickers[i] = ticker.Ticker
	}
	return tickers, nil
}

// Delete a ticker from user's watchlist
// Require an authenticated user to perform this action
func DeleteTickerFromWatchlist(ticker string) error {
	// Ensure the current user is authenticated
	resp, _ := Client.Auth.GetUser()
	if resp == nil {
		return errors.New("User not logged in")
	}

	// Delete the relevant entry from watchlists table
	_, count, err := Client.From("watchlists").
		Delete("", "exact").
		Eq("user_id", resp.User.ID.String()).
		Eq("ticker", ticker).
		Execute()
	if err != nil {
		return err
	} else if count == 0 {
		return errors.New("Ticker not in user's watchlist")
	}
	return err
}
