package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"strconv"
	"fmt"

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
// @Router /events/{ticker} [get]
func GetEventsHandler(c *gin.Context) {
	ticker := c.Param("ticker")
	fromStr := c.Param("from")
	toStr := c.Param("to")
	limitStr := c.Param("limit")

	if ticker == "" || fromStr == "" || toStr == "" || limitStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required query parameters"})
        return
    }

    from, err := time.Parse(time.RFC3339, fromStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' time format"})
        return
    }
    to, err := time.Parse(time.RFC3339, toStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' time format"})
        return
    }
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'limit' parameter"})
        return
    }

	events, _ := supabase.GetEvents(ticker, from, to, limit)

	// Serialize (marshal) to JSON
	jsonData, err := json.Marshal(events)
	if err != nil {
		fmt.Println("Error serializing:", err)
		return
	}
	c.JSON(http.StatusOK, jsonData)
}