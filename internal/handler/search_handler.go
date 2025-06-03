package handler

import (
	// "log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	// "omnichart-server/internal/supabase"
	// "omnichart-server/internal/integration/marketdata"
)

// Mock database (replace with actual DB/logic)
var tickers = []map[string]string{
	{"ticker": "NVDA", "name": "NVIDIA Corporation"},
	{"ticker": "TSM", "name": "Taiwan Semiconductor Manufacturing Company Limited"},
	{"ticker": "SMH", "name": "VanEck Semiconductor ETF"},
	{"ticker": "AMD", "name": "Advanced Micro Devices"},
}

// SearchHandler godoc
// @Summary      Search for tickers by name or symbol
// @Description  Returns a list of matching ticker objects whose name or symbol contains the query string
// @Tags         search
// @Accept       json
// @Produce      json
// @Param        q     query     string  true  "Search query (e.g., part of a ticker symbol or company name)"
// @Success      200   {object}  map[string][]map[string]string  "List of matching ticker entries"
// @Failure      400   {object}  map[string]string               "Missing or invalid query parameter"
// @Router       /search [get]
func GetSearchHandler(c *gin.Context) {
	query := strings.ToUpper(c.Query("q"))
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameter `q`"})
		return
	}

	var results []map[string]string

	// asset, err := marketdata.Client.GetAsset(query)
	// if err != nil {
	// 	log.Println(err, asset)
	// }

	// Filter tickers

	for _, ticker := range tickers {
		if strings.Contains(strings.ToLower(ticker["ticker"]), query) || strings.Contains(strings.ToLower(ticker["name"]), query) {
			results = append(results, ticker)
		}
	}

	c.JSON(http.StatusOK, results)
}
