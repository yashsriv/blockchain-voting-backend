package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	radix "github.com/mediocregopher/radix/v3"

	"blockchain-voting/redis"
)

type registerRequest struct {
	PrivateKey string `json:"public"`
	PublicKey  string `json:"private"`
}

func Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request registerRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var pubKeyName, privKeyName string
		isAdmin := ctx.GetBool("isAdmin")
		isCandidate := ctx.GetBool("isCandidate")

		if !isAdmin && !isCandidate {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "neither admin nor candidate",
			})
			return
		}

		if isAdmin {
			pubKeyName = "admin-pub"
			privKeyName = "admin-priv"
		} else {
			username := ctx.GetString("username")
			pubKeyName = fmt.Sprintf("candidate-pub-%s", username)
			privKeyName = fmt.Sprintf("candidate-priv-%s", username)
		}

		err := redis.Client.Do(radix.Cmd(nil, "SET", pubKeyName, request.PublicKey))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = redis.Client.Do(radix.Cmd(nil, "SET", privKeyName, request.PrivateKey))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// TODO: Perform solidity transaction to publish pubkey and return link
		// of transaction to user
		ctx.JSON(http.StatusOK, gin.H{
			"link": "todo://ethereum",
		})
	}
}
