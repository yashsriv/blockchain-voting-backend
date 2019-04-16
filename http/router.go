package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"blockchain-voting/http/handlers"
	"ethlib"
)

// Router is the root-level router used by the
// server
var VC *ethlib.VotingContractWrapper

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
	router.GET("/all-votes", handlers.AuthCheck(), handlers.GetAllVotes(VC))
	router.GET("/all-voters", handlers.AuthCheck(), handlers.GetAllVoters(VC))
	router.POST("/end-voting", handlers.AuthCheck(), handlers.EndVoting(VC))
	router.POST("/start-voting", handlers.AuthCheck(), handlers.StartVoting(VC))
	router.POST("/vote", handlers.AuthCheck(), handlers.Vote(VC))
	router.POST("/publish-results", handlers.AuthCheck(), handlers.PublishResult(VC))

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
