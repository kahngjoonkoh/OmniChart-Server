package supabase

import (
    "log"
    "os"

	"github.com/joho/godotenv"
    "github.com/supabase-community/supabase-go"
)

var Client *supabase.Client

func Init() {
	var err error
	err = godotenv.Load()

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
