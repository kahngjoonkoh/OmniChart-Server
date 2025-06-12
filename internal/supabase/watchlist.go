package supabase

import (
	"encoding/json"
	"fmt"

	"omnichart-server/internal/models"
)

// Add a ticker to user's watchlist
// Require an authenticated user to perform this action
func AddTickerToWatchlist(token, ticker string) error {
	// Ensure the current user is authenticated
	client := Client.Auth.WithToken(token)
	resp, err := client.GetUser()
	if err != nil {
		return fmt.Errorf("failed to authenticate user")
	}

	// Ensure there is no existing entry
	inWatchlist, err := TickerInWatchlist(token, ticker)
	if err != nil {
		return fmt.Errorf("failed to fetch watchlist")
	} else if inWatchlist {
		return fmt.Errorf("%s already in watchlist", ticker)
	}

	// Insert new entry to watchlists table
	req := map[string]interface{}{
		"user_id"	: resp.User.ID.String(),
		"ticker"	: ticker,
	}
	_, _, err = Client.From("watchlists").
		Insert(req, false, "ignore", "", "").
		Execute()
	if err != nil {
		return fmt.Errorf("failed to add %s to watchlist", ticker)
	}
	return nil
}

// Retrieve the tickers in user's watchlist
// Require an authenticated user to perform this action
func GetWatchlist(token string) ([]string, error) {
	// Fetch the authenticated user
	client := Client.Auth.WithToken(token)
	userResp, err := client.GetUser()
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to authenticate user")
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

// Remove a ticker from user's watchlist
// Require an authenticated user to perform this action
func RemoveTickerFromWatchlist(token, ticker string) error {
	// Ensure the current user is authenticated
	client := Client.Auth.WithToken(token)
	resp, err := client.GetUser()
	if err != nil {
		return fmt.Errorf("failed to authenticate user")
	}

	// Delete the relevant entry from watchlists table
	_, count, err := Client.From("watchlists").
		Delete("", "exact").
		Eq("user_id", resp.User.ID.String()).
		Eq("ticker", ticker).
		Execute()
	if err != nil {
		return fmt.Errorf("failed to remove %s from watchlist", ticker)
	} else if count == 0 {
		return fmt.Errorf("%s not in user's watchlist", ticker)
	}
	return nil
}

// Check whether a ticker is in user's watchlist
// Require an authenticated user to perform this action
func TickerInWatchlist(token, ticker string) (bool, error) {
	// Ensure the current user is authenticated
	client := Client.Auth.WithToken(token)
	resp, err := client.GetUser()
	if err != nil {
		return false, fmt.Errorf("failed to authenticate user")
	}

	_, count, err := Client.From("watchlists").
		Select("", "exact", false).
		Eq("user_id", resp.User.ID.String()).
		Eq("ticker", ticker).
		Execute()
	if err != nil {
		return false, fmt.Errorf("failed to fetch user's watchlist")
	}
	return count > 0, nil
}
