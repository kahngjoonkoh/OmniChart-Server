package main

import (
   "os"
   "log"

	"github.com/joho/godotenv"

   "omnichart-server/internal/router"
   "omnichart-server/internal/supabase"
   "omnichart-server/internal/alpaca"
)


func main() {
   err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load; falling back to environment variables")
	}
   supabase.Init()
   alpacaApi.Init()

   r := router.SetupRouter()

   port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // default for local dev
    }
   r.Run(":" + port)
}