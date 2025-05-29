package main

import (
   "os"

   "omnichart-server/internal/router"
   "omnichart-server/internal/supabase"
)


func main() {
   supabase.Init()

   r := router.SetupRouter()

   port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // default for local dev
    }
   r.Run(":" + port)
}