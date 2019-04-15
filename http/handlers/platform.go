package handlers

import (
	"fmt"
	"net/http"

	"blockchain-voting/redis"
	"github.com/gin-gonic/gin"
	radix "github.com/mediocregopher/radix/v3"
)

func GetInfo() gin.HandlerFunc {
	func(ctx *gin.Context) {
		// Get admin pubKey from database
		var adminKey string
		err := client.Do(radix.Cmd(&adminKey, "GET", "admin-pub"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Get candidate list
		candidates := viper.GetStringSlice("candidates")

		// Populate map with pubKey of candidates
		var candidateKeys = make(map[string]string)
		for _, candidate := range candidates {
			candidateMapKey := fmt.Sprintf("candidate-pub-%s", candidate)
			var candidatePubKey string
			err := client.Do(radix.Cmd(&candidatePubKey, "GET", candidateMapKey))
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			candidateKeys[candidate] = candidatePubKey
		}

		// Get if voting has started
		var votingStarted bool
		err := client.Do(radix.Cmd(&votingStarted, "GET", "votingStarted"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Get if voting has ended
		var votingEnded bool
		err := client.Do(radix.Cmd(&votingEnded, "GET", "votingStarted"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Get if results have been published
		var resultsPublished bool
		err := client.Do(radix.Cmd(&resultsPublished, "GET", "votingStarted"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H {
			"adminKey": adminKey,
			"candidateKeys": candidateKeys,
			"candidates": candidates,
			"votingStarted": votingStarted,
			"votingEnded": votingEnded,
			"resultsPublished": resultsPublished,
		})
	}
}
