package supabase

import (
	"strings"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/supabase-community/gotrue-go/types"
)

type SupabaseError struct {
	Message	  string	`json:"message"`
}

// Create a new user and corresponding profile given username, email and password
func SignUpUser(username, email, password string) error {
	req := types.SignupRequest{
		Email: email,
		Password: password,
		Data: map[string]interface{}{
			"username": username,
		},
	}
	_, err := Client.Auth.Signup(req)
	if err != nil {
		// Retrieve message from error
		msg := strings.ReplaceAll(err.Error(), "msg", "message")
		fmt.Println(msg)
		index := strings.Index(msg, "{")
		var supaErr SupabaseError
		// Return raw error if json parsing fails
		if unmarshalErr := json.Unmarshal([]byte(msg[index:]), &supaErr); unmarshalErr != nil {
			return err
		}
		msg = supaErr.Message

		// Handle specific error messages
		if strings.Contains(msg, "User already registered") {
			return errors.New("Email already registered.")
		} else if strings.Contains(msg, "duplicate key value violates unique constraint") {
			return errors.New("Username already registered.")
		} else {
			return errors.New(msg)
		}
	}
	return nil
}