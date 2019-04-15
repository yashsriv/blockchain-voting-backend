package handlers

import (
	"blockchain-voting/redis"
	"net/http"

	"github.com/gin-gonic/gin"
	radix "github.com/mediocregopher/radix/v3"
)

const votersList = "votersList"
const votes = "votes"

type VoteRequest struct {
	Vote string `json:"vote"`
}

func Vote() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request VoteRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err := redis.Client.Do(radix.Cmd(nil, "RPUSH", request.Vote))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		var votingStarted bool
		err = redis.Client.Do(radix.Cmd(&votingStarted, "GET", "votingStarted"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		var votingEnded bool
		err = redis.Client.Do(radix.Cmd(&votingEnded, "GET", "votingEnded"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Check if voting has been started or not
		if !(votingStarted && !votingEnded) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "voting has not been started yet",
			})
			return
		}

		// Checking if the user has voted or not
		var voted int
		err = redis.Client.Do(radix.Cmd(&voted, "SISMEMBER", votersList))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if voted == 1 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}

		// if not voted
		err = redis.Client.Do(radix.Cmd(nil, "SADD", votersList, ctx.GetString("username")))
		// TODO: Perform solidity transaction to vote and publish the transaction to the user
		ctx.JSON(http.StatusOK, gin.H{
			"link": "todo://ethereum",
		})

	}
}
