package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TagData(tag string) string {
	return tag
}

// Path /tags/[tag]
func GetStocksByTag(c *gin.Context) {
	c.Param("tags")
	tag := c.Param("tag")
	info := TagData(tag)
	c.IndentedJSON(http.StatusOK, info)
}

// Path /tags
func GetTag(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Tags landing page")
}
