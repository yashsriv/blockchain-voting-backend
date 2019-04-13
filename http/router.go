package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"blockchain-voting/http/handlers"
)

// Router is the root-level router used by the
// server
func Router() *gin.Engine {
	// Logger and Recovery Middleware already loaded
	router := gin.Default()

	router.POST("/login", handlers.HandleLogin())
	router.GET("/ping", handlers.AuthCheck(), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"pong": "Hello World",
		})
	})

	return router
}
