package router

import (
	"github.com/gin-gonic/gin"

	docs "omnichart-server/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"omnichart-server/internal/handler"
)

// @BasePath /api/v1
func SetupRouter() *gin.Engine {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"

	v1 := r.Group("/api/v1")
	v1.GET("/events/:ticker", handler.GetEventsHandler)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
