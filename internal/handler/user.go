package handler

import (
	"net/http"

	"omnichart-server/internal/models"
	"omnichart-server/internal/supabase"

	"github.com/gin-gonic/gin"
)

// @Summary Sign up a new user account
// @Accept json
// @Produce json
// @Param signup body object true "Signup data"
// @Success 200 {array} models.UserInfo
// @Failure 400 {object} map[string]interface{}
// @Router /signup [post]
func SignUpHandler(c *gin.Context) {
	var req models.SignUpRequest

	// Extract user data from request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Create new user in supabase
	err := supabase.SignUpUser(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return user info
	c.JSON(http.StatusCreated, models.UserInfo{
		Username: req.Username,
	})
}