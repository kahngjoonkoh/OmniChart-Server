package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"omnichart-server/internal/supabase"
)

// PostTickerEventRequest defines the JSON body for POST /api/v1/ticker_events
type PostTickerEventRequest struct {
	Ticker     string `json:"ticker" example:"AAPL"`
	EventID    string `json:"event_id" example:"e7b8c6e5-3e44-4e2a-a6e7-123456789abc"`
	StartIndex int    `json:"start_index" example:"0"`
	EndIndex   int    `json:"end_index" example:"50"`
}

// @Summary      Add a new ticker event segment
// @Description  Inserts a ticker event segment for a given ticker and event.
// @Tags         ticker_events
// @Accept       json
// @Produce      json
// @Param        tickerEvent  body      PostTickerEventRequest  true  "Ticker Event payload"
// @Success      200          {object}  models.TickerEvent
// @Failure      400          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Router       /api/v1/ticker_events [post]
func PostTickerEventHandler(c *gin.Context) {
	var req PostTickerEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	tickerEvent, err := supabase.AddTickerEvent(req.Ticker, req.EventID, req.StartIndex, req.EndIndex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to post ticker event"})
		return
	}

	c.JSON(http.StatusOK, tickerEvent)
}

// @Summary      List ticker events for a ticker
// @Description  Returns all ticker event segments associated with a ticker symbol.
// @Tags         ticker_events
// @Accept       json
// @Produce      json
// @Param        ticker  path      string  true  "Ticker symbol"
// @Success      200     {array}   models.TickerEvent
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /api/v1/ticker_events/{ticker} [get]
func GetTickerEventsHandler(c *gin.Context) {
	ticker := c.Param("ticker")
	if ticker == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing ticker"})
		return
	}

	tickerEvents, err := supabase.GetTickerEvents(ticker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch ticker events"})
		return
	}

	c.JSON(http.StatusOK, tickerEvents)
}
