package handlers

import (
	"blockchain-voting/redis"
	"net/http"

	"github.com/gin-gonic/gin"
	radix "github.com/mediocregopher/radix/v3"
	"github.com/spf13/viper"
)

const VotersList = "votersList"
const Votes = "votes"

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

		err := redis.Client.Do(radix.Cmd(nil, "RPUSH", Votes, request.Vote))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		var votingStarted bool
		err = redis.Client.Do(radix.Cmd(&votingStarted, "GET", VotingStarted))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		var votingEnded bool
		err = redis.Client.Do(radix.Cmd(&votingEnded, "GET", VotingEnded))
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
		err = redis.Client.Do(radix.Cmd(&voted, "SISMEMBER", VotersList))
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
		err = redis.Client.Do(radix.Cmd(nil, "SADD", VotersList, ctx.GetString("username")))
		// TODO: Perform solidity transaction to vote and publish the transaction to the user
		ctx.JSON(http.StatusOK, gin.H{
			"link": "todo://ethereum",
		})

	}
}

func StartVoting() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.GetString("username")
		admin := viper.GetString("admin.username")

		if admin != username {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "you are not admin",
			})
			return
		}

		err := redis.Client.Do(radix.Cmd(nil, "SET", VotingStarted, "true"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		// TODO : perform the solidity transsactions to start the voting

		ctx.JSON(http.StatusOK, gin.H{
			"link": "todo://ethereum",
		})

	}
}

func EndVoting() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.GetString("username")
		admin := viper.GetString("admin.username")

		if admin != username {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "you are not admin",
			})
			return
		}

		err := redis.Client.Do(radix.Cmd(nil, "SET", VotingEnded, "true"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		// TODO : perform the solidity transsactions to start the voting

		ctx.JSON(http.StatusOK, gin.H{
			"link": "todo://ethereum",
		})

	}
}

func GetAllVotes() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var votes []string
		err := redis.Client.Do(radix.Cmd(&votes, "LRANGE", Votes, "0", "-1"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		ctx.JSON(http.StatusOK, votes)
	}
}
