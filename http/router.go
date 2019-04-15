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

	// Login
	router.POST("/login", handlers.HandleLogin())

	// Register admins and candidates
	router.POST("/register", handlers.AuthCheck(), handlers.Register())

	// Get Platform info
	router.GET("/platform-info", handlers.AuthCheck(), handlers.GetInfo())

	// Voting handlers
	router.GET("/get-all-votes", handlers.AuthCheck(), handlers.GetAllVotes())
	router.POST("/end-voting", handlers.AuthCheck(), handlers.EndVoting())
	router.POST("/start-voting", handlers.AuthCheck(), handlers.StartVoting())
	router.POST("/vote", handlers.AuthCheck(), handlers.Vote())

	// Get encrypted-admin-privKey
	router.GET("/admin-privKey", handlers.AuthCheck(), handlers.GetAdminPrivKey())
	// Get encrypted-(particular)candidate-privKey
	router.GET("/candidate-privKey", handlers.AuthCheck(), handlers.GetCandidatePrivKey())

	router.GET("/ping", handlers.AuthCheck(), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"pong": ctx.GetString("username"),
		})
	})

	return router
}
