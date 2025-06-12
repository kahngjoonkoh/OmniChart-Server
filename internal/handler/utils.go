package handler

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// Get the user given the access token from the authorization header of a request
func GetAccessToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", fmt.Errorf("invalid authorization header %s", authHeader)
	}
	return parts[1], nil
}
