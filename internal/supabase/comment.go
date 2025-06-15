package supabase

import (
	"encoding/json"
	"fmt"
	"omnichart-server/internal/models"

	"github.com/supabase-community/postgrest-go"
)

// Add a comment entry to database
// Require an authenticated user to perform this action
func AddComment(token, tickerEventID, content, sentiment string) (*models.Comment, error) {
	// Ensure the current user is authenticated
	client := Client.Auth.WithToken(token)
	userResp, err := client.GetUser()
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate user")
	}

	// Fetch the username
	username, err := GetUsername(userResp.User.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve username")
	}

	// Create insertion request
	req := map[string]interface{}{
		"content":         content,
		"ticker_event_id": tickerEventID,
		"user_id":         userResp.User.ID.String(),
		"username":		   username,
		"sentiment":	   sentiment,
	}
	resp, _, err := Client.From("comments").
		Insert(req, false, "representation", "", "").
		Single().
		Execute()
	if err != nil {
		return nil, err
	}

	// Extract data from json response
	var comment models.Comment
	err = json.Unmarshal(resp, &comment)
	if err != nil {
		return nil, err
	}
	fmt.Println(comment)

	return &comment, nil
}

// Retrieve the comments of a ticker event sorted by time
func GetComments(tickerEventID, filter string, ascending bool) ([]models.Comment, error) {
	// Compose order options for ascending order
	orderOpts := &postgrest.OrderOpts{Ascending: ascending}

	// Select params: columns, head ("" for no head), count (false)
	query := Client.From("comments").
		Select("*", "exact", false).
		Eq("ticker_event_id", tickerEventID)
	if filter != "" {
		query = query.Eq("sentiment", filter)
	}
	resp, count, err := query.
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