package handler

import (
	"net/http"
	"strings"
	"time"

	"omnichart-server/internal/alpaca"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/gin-gonic/gin"
)

func GetHistoricalDataHandler(c *gin.Context) {
	ticker := strings.ToUpper(c.Param("ticker"))

	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'start' or 'end' query parameters. Use RFC3339 format."})
		return
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'start' format. Use RFC3339"})
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'end' format. Use RFC3339"})
		return
	}

	if end.Before(start) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'end' must be after 'start'"})
		return
	}

	// ✅ Correct API usage for v3
	data, err := alpacaApi.MarketData.GetBars(ticker, marketdata.GetBarsRequest{
		TimeFrame: marketdata.OneDay,
		Start:     start,
		End:       end,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching historical bars", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data) // ✅ Return slice directly
}
