package handler

import (
	"net/http"

	"omnichart-server/internal/models"
	"omnichart-server/internal/supabase"

	"github.com/gin-gonic/gin"
)

// @Summary Add a new ticker to current user's watchlist
// @Accept json
// @Produce json
// @Param ticker body object true "Ticker"
// @Param token header string false "Access token"
// @Router /watchlist/add [post]
func AddTickerHandler(c *gin.Context) {
	var req models.Watchlist

	// Extract data from user request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call methods to add new entry in database
	if err := supabase.AddTickerToWatchlist(req.Ticker); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// @Summary Retrieve the tickers from current user's watchlist
// @Param token header string false "Access token"
// @Router /watchlist [get]
func GetWatchlistHandler(c *gin.Context) {
	tickers, err := supabase.GetWatchlist()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tickers": tickers,
	})
}

// @Summary Add a new ticker to current user's watchlist
// @Accept json
// @Produce json
// @Param ticker body object true "Ticker"
// @Param token header string false "Access token"
// @Router /watchlist/remove [delete]
func RemoveTickerHandler(c *gin.Context) {
	var req models.Watchlist

	// Extract ticker code from user's request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call methods to delete relevant table entry
	if err := supabase.DeleteTickerFromWatchlist(req.Ticker); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}
