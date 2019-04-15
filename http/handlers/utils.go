package handlers

import (
	"net/http"

	"blockchain-voting/redis"
	"github.com/gin-gonic/gin"
	radix "github.com/mediocregopher/radix/v3"
)

func GetAdminPrivKey() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get admin privKey from database
		var adminPrivKey string
		err := redis.Client.Do(radix.Cmd(&adminPrivKey, "GET", "admin-priv"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"encrypted-admin-privKey": adminPrivKey,
		})
	}
}
