package main

import (
   "omnichart-server/internal/router"
   "omnichart-server/internal/supabase"
)


func main() {
   supabase.Init()

   r := router.SetupRouter()
   r.Run(":8080")
}