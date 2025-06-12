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
	// Extract access token from header
	token, err := GetAccessToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid authorization token"})
		return
	}

	// Extract data from user request
	var req models.Watchlist
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call methods to add new entry in database
	if err := supabase.AddTickerToWatchlist(token, req.Ticker); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// @Summary Retrieve the tickers from current user's watchlist
// @Param token header string false "Access token"
// @Router /watchlist [get]
func GetWatchlistHandler(c *gin.Context) {
	// Extract access token from header
	token, err := GetAccessToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid authorization token"})
		return
	}

	tickers, err := supabase.GetWatchlist(token)
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
	// Extract access token from header
	token, err := GetAccessToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid authorization token"})
		return
	}

	// Extract ticker code from user's request
	var req models.Watchlist
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call methods to delete relevant table entry
	if err := supabase.RemoveTickerFromWatchlist(token, req.Ticker); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// @Summary Check if a ticker is in user's watchlist
// @Param ticker path string true "Ticker symbol"
// @Param token header string false "Access token"
// @Router /watchlist/{ticker} [get]
func TickerInWatchlistHandler(c *gin.Context) {
	// Extract access token from header
	token, err := GetAccessToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid authorization token"})
		return
	}

	ticker := c.Param("ticker")
	inWatchlist, err := supabase.TickerInWatchlist(token, ticker)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"in": inWatchlist})
}
