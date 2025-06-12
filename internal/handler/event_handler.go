package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"omnichart-server/internal/supabase"
)

// GetEventsHandler godoc
// @Summary Gets events for a given ticker and timeframe
// @Param ticker path string true "Ticker symbol"
// @Param from query string true "Start time in RFC3339 format"
// @Param to query string true "End time in RFC3339 format"
// @Param limit query int true "Maximum number of events to return"
// @Success 200 {array} models.Event
// @Failure 400 {object} map[string]interface{}
// @Router /events/{event_id} [get]
func GetEventsHandler(c *gin.Context) {
	ticker := strings.ToUpper(c.Param("ticker"))

	events, err := supabase.GetEvents(ticker)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Fetching from database"})
		return
	}

	c.JSON(http.StatusOK, events)
}
