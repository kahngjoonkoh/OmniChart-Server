package router

import (
    "net/http"

    "github.com/gin-gonic/gin"

    docs "omnichart-server/docs"
    ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
)

// @BasePath /api/v1

// optional: description, tags, accept, produce.
// Helloworld godoc
// @Summary ping example
// @Success 200 {string} string "helloworld"
// @Router /example/helloworld [get]
func Helloworld(c *gin.Context) {
    c.JSON(http.StatusOK, "helloworld")
}

func SetupRouter() *gin.Engine {
    r := gin.Default()

    docs.SwaggerInfo.BasePath = "/api/v1"

    v1 := r.Group("/api/v1")
    {
        eg := v1.Group("/example")
        eg.GET("/helloworld", Helloworld)
    }

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    return r
}