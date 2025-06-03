package router

import (
	"github.com/gin-gonic/gin"

	docs "omnichart-server/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"omnichart-server/internal/handler"
	"omnichart-server/internal/models"
)

// @BasePath /api/v1
func SetupRouter() *gin.Engine {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"

	v1 := r.Group("/api/v1")
	v1.GET("/events/:ticker", handler.GetEventsHandler)
	v1.GET("/tags", models.GetTag)              // Landing page for all tags
	v1.GET("/tags/:tag", models.GetStocksByTag) // Page for specific tag
	v1.POST("/comments", handler.PostCommentHandler)
    v1.GET("/comments/:tickerEventID", handler.GetCommentsHandler)
	v1.GET("/search", handler.GetSearchHandler)
	v1.POST("/ticker_events", handler.PostTickerEventHandler)
	v1.GET("/ticker_events/:ticker", handler.GetTickerEventsHandler)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
