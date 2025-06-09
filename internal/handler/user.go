package handler

import (
	"net/http"
	"strings"

	"omnichart-server/internal/models"
	"omnichart-server/internal/supabase"

	"github.com/gin-gonic/gin"
)

// @Summary Sign up a new user account
// @Accept json
// @Produce json
// @Param signup body object true "Username, email and password"
// @Success 200 {array} models.UserInfo
// @Failure 400 {object} map[string]interface{}
// @Router /signup [post]
func SignUpHandler(c *gin.Context) {
	var req models.SignUpRequest

	// Extract user data from request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

// @Summary Login a user
// @Accept json
// @Produce json
// @Param login body object true "Username and password"
// @Success 200 {array} models.UserInfo
// @Failure 400 {object} map[string]interface{}
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Tries to login user with provided username and password
	accessToken, refreshToken, err := supabase.LoginUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.UserInfo{
		Username: req.Username,
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	})
}

// @Summary Log out a user
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /logout [post]
func LogoutHandler(c *gin.Context) {
	err := supabase.LogoutUser()
	if err != nil {
		msg := err.Error()
		if !strings.Contains(msg, "This endpoint requires a Bearer token") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, map[string]interface{}{})
}
