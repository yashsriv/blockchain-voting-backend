package handlers

import (
	"fmt"
	"net/http"

	"blockchain-voting/redis"
	"github.com/gin-gonic/gin"
	radix "github.com/mediocregopher/radix/v3"
)

const AdminPriv = "admin-priv"
const CandidatePriv = "candidate-priv"

func GetAdminPrivKey() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get admin privKey from database
		var adminPrivKey string
		err := redis.Client.Do(radix.Cmd(&adminPrivKey, "GET", AdminPriv))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		ctx.JSON(http.StatusOK, adminPrivKey)
	}
}

func GetCandidatePrivKey() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get username from the Context
		username := ctx.Query("username")

		// Construct redis key
		candidateRedisKey := fmt.Sprintf("%s-%s", CandidatePriv, username)

		// Get candidate privKey from Database
		var candPrivKey string
		err := redis.Client.Do(radix.Cmd(&candPrivKey, "GET", candidateRedisKey))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		ctx.JSON(http.StatusOK, candPrivKey)
	}
}
