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

type EmailResponse struct {
	Email	string	`json:"email"`
}

// Login a user given username and password
func LoginUser(username, password string) (string, string, error) {
	// Find the email from profiles table
	resp, _, err := Client.From("profiles").
		Select("email", "", false).
		Eq("username", username).
		Single().
		Execute()
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "rows returned") {
			return "", "", errors.New("Invalid username or password")
		}
		return "", "", err
	}
	
	// Retrieve email from the response
	var email EmailResponse
	err = json.Unmarshal(resp, &email)
	if err != nil {
		return "", "", err
	}

	// Login the user
	session, err := Client.SignInWithEmailPassword(email.Email, password)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "Invalid login credentials") {
			return "", "", errors.New("Invalid username or password")
		}
		fmt.Println(err.Error())
		return "", "", err
	}
	return session.AccessToken, session.RefreshToken, nil
}

// Logout a user
func LogoutUser() error {
	return Client.Auth.Logout()
}
