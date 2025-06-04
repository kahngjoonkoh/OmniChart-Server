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
	// timeframe := c.Param("timeframe")

	data, err := alpacaApi.MarketData.GetBars(ticker, marketdata.GetBarsRequest{
		TimeFrame: marketdata.OneDay,
		Start: time.Now().AddDate(0, -6, 0),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Fetching Historical Bars."})
		return
	}
	
	c.JSON(http.StatusOK, data)
}

func GetLiveDataHandler(c *gin.Context) {

}
