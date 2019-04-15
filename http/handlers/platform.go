package handlers

import (
	"fmt"
	"net/http"

	"blockchain-voting/redis"

	"github.com/gin-gonic/gin"
	radix "github.com/mediocregopher/radix/v3"
	"github.com/spf13/viper"
)

const VotingStarted = "votingStarted"
const VotingEnded = "votingEnded"
const IsResultPublished = "resultPublished"

func GetInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get admin pubKey from database
		var adminKey string
		err := redis.Client.Do(radix.Cmd(&adminKey, "GET", "admin-pub"))
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
			err = redis.Client.Do(radix.Cmd(&candidatePubKey, "GET", candidateMapKey))
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
		err = redis.Client.Do(radix.Cmd(&votingStarted, "GET", VotingStarted))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Get if voting has ended
		var votingEnded bool
		err = redis.Client.Do(radix.Cmd(&votingEnded, "GET", VotingEnded))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Get if results have been published
		var resultsPublished bool
		err = redis.Client.Do(radix.Cmd(&resultsPublished, "GET", IsResultPublished))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"adminKey":         adminKey,
			"candidateKeys":    candidateKeys,
			"candidates":       candidates,
			"votingStarted":    votingStarted,
			"votingEnded":      votingEnded,
			"resultsPublished": resultsPublished,
		})
	}
}
