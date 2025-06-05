package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"fmt"

	"github.com/gin-gonic/gin"

	alpacaApi "omnichart-server/internal/alpaca"
	"omnichart-server/internal/models"
	"omnichart-server/internal/supabase"
)

// GetSearchHandler godoc
// @Summary      Search for tickers by symbol or company name
// @Description  Returns matching ticker objects whose symbol or company name contains the query string.
// @Tags         search
// @Accept       json
// @Produce      json
// @Param        q     query     string  true  "Search query string (part of ticker symbol or company name)"
// @Success      200   {object}  map[string][]models.Ticker  "Search results grouped by type (e.g., stocks)"
// @Failure      400   {object}  map[string]string           "Missing or invalid query parameter"
// @Failure      500   {object}  map[string]string           "Internal server error (DB or API failure)"
// @Router       /search [get]
func GetSearchHandler(c *gin.Context) {
	query := strings.ToUpper(c.Query("q"))
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameter `q`"})
		return
	}

	data, _, err := supabase.Client.From("tickers").
		Select("*", "", false).
		Or(fmt.Sprintf("ticker.ilike.%s*,name.ilike.*%s*", query, query), "").
		Execute()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch from DB"})
		return
	}

	var tickers []models.Ticker
	if err := json.Unmarshal(data, &tickers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unmarshal Error. Received incorrect datatype"})
		return
	}

	if len(tickers) == 0 {
		log.Println("No ticker found. Searching Alpaca API")
		asset, err := alpacaApi.Client.GetAsset(query)
		if err != nil {
			log.Println(err, asset)
		}
		log.Println(asset)

		newTicker := models.Ticker{
			Ticker: asset.Symbol,
			Name:   asset.Name,
		}

		if _, _, err := supabase.Client.From("tickers").Insert(newTicker, false, "representation", "", "").Execute(); err != nil {
			log.Println("Error inserting ticker:", err)
			return
		}

		log.Println("Ticker inserted successfully")

		tickers = append(tickers, newTicker)
	}

	filteredStocks := make([]models.Ticker, 0)


	has_ticker_result := false

	for _, ticker := range tickers {
		if strings.HasPrefix(ticker.Ticker, query) {
			filteredStocks = append(filteredStocks, ticker)
			has_ticker_result = true
		}
	}

	for _, ticker := range tickers {
		if !has_ticker_result && strings.Contains(strings.ToUpper(ticker.Name), query) {
			filteredStocks = append(filteredStocks, ticker)
		}
	}

	results := map[string]interface{}{
		"stocks": filteredStocks,
		// "events": filteredEvents,
		// "others": filteredOthers,
	}

	c.JSON(http.StatusOK, results)
}
