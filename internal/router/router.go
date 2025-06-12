package router

import (
	"time"

	"github.com/gin-contrib/cors"
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173", "https://omnichart.impaas.uk"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // required if using cookies or authorization headers
		MaxAge:           12 * time.Hour,
	}))

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
	v1.GET("/bars/:ticker", handler.GetHistoricalDataHandler)
	// v1.GET("/bars/:ticker", handler.GetLiveDataHandler)
	v1.POST("/signup", handler.SignUpHandler) // User sign up
	v1.POST("/login", handler.LoginHandler) // User login
	v1.POST("/logout", handler.LogoutHandler) // User logout#

	v1.GET("/beta/:ticker", handler.GetBetaHandler)
	v1.POST("/watchlist/add", handler.AddTickerHandler) // Add ticker to watchlist
	v1.GET("/watchlist", handler.GetWatchlistHandler) // Get tickers in watchlist
	v1.GET("/watchlist/:ticker", handler.TickerInWatchlistHandler)
	v1.DELETE("/watchlist/remove", handler.RemoveTickerHandler) // Remove ticker from watchlist

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
