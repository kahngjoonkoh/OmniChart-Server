package main

import (
   "os"

   "omnichart-server/internal/router"
   "omnichart-server/internal/integration/db"
   "omnichart-server/internal/integration/marketdata"
)


func main() {
   db.Init()
   marketdata.Init()

   r := router.SetupRouter()

   port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // default for local dev
    }
   r.Run(":" + port)
}