// internal/handler/comments.go

package handler

import (
	"net/http"
	"omnichart-server/internal/supabase"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PostCommentRequest defines the JSON body for POST /api/v1/comments
type PostCommentRequest struct {
	TickerEventID string `json:"ticker_event_id" example:"drop1"`
	Content       string `json:"content" example:"This drop makes sense after Q2 miss."`
	Sentiment	  string `json:"sentiment"`
}

// @Summary      Add a new comment
// @Description  Inserts a comment for the given ticker event segment.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        comment  body      PostCommentRequest  true  "Comment payload"
// @Param		 Authorization	  header	string		false "Authorization header"
// @Success      200      {object}  models.Comment
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /comments [post]
func PostCommentHandler(c *gin.Context) {
	// Extract access token from header
	token, err := GetAccessToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid authorization token"})
		return
	}

	var req PostCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	comment, err := supabase.AddComment(token, req.TickerEventID, req.Content, req.Sentiment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to post comment",
			"details": err.Error(), // optional, remove in production
		})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// @Summary      List comments for an event
// @Description  Returns all comments associated with a ticker event ID.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        tickerEventID  path      string  true  "Ticker Event ID"
// @Success      200            {array}   models.Comment
// @Failure      400            {object}  map[string]string
// @Failure      500            {object}  map[string]string
// @Router       /comments/{tickerEventID} [get]
func GetCommentsHandler(c *gin.Context) {
	tickerEventID := c.Param("tickerEventID")
	if tickerEventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tickerEventID"})
		return
	}

	sentiment := c.Query("sentiment")
	ascending, err := strconv.ParseBool(c.Query("ascending"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ascending must be either true or false"})
		return
	}

	comments, err := supabase.GetComments(tickerEventID, sentiment, ascending)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}
