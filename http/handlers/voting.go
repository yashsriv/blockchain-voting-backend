package handlers

import (
	"blockchain-voting/redis"
	"net/http"
	"math/big"

	"ethlib"

	"github.com/gin-gonic/gin"
	radix "github.com/mediocregopher/radix/v3"
	"github.com/spf13/viper"
)

const VotersList = "votersList"
const Votes = "votes"

type VoteRequest struct {
	Vote string `json:"vote"`
}

func Vote(vc *ethlib.VotingContractWrapper) gin.HandlerFunc {
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

		var votingStarted string
		err = redis.Client.Do(radix.Cmd(&votingStarted, "GET", VotingStarted))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		var votingEnded string
		err = redis.Client.Do(radix.Cmd(&votingEnded, "GET", VotingEnded))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Check if voting has been started or not
		if !(votingStarted == "1" && votingEnded == "0") {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "voting has not started yet",
			})
			return
		}

		// Checking if the user has voted or not
		var voted int
		err = redis.Client.Do(radix.Cmd(&voted, "SISMEMBER", VotersList, ctx.GetString("username")))
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

		// Solidity vote interaction
		txhash, err := vc.AddEncryptedVote(request.Vote, ctx.GetString("username"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// if not voted
		err = redis.Client.Do(radix.Cmd(nil, "SADD", VotersList, ctx.GetString("username")))

		ctx.JSON(http.StatusOK, gin.H{
			"link": txhash,
		})

	}
}

func StartVoting(vc *ethlib.VotingContractWrapper) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.GetString("username")
		admin := viper.GetString("admin.username")

		if admin != username {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "you are not admin",
			})
			return
		}

		err := redis.Client.Do(radix.Cmd(nil, "SET", VotingStarted, "1"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = redis.Client.Do(radix.Cmd(nil, "SET", VotingEnded, "0"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// TODO : perform the solidity transsactions to start the voting
		txhash, err := vc.StartVoting()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"link": txhash,
		})

	}
}

func EndVoting(vc *ethlib.VotingContractWrapper) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.GetString("username")
		admin := viper.GetString("admin.username")

		if admin != username {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "you are not admin",
			})
			return
		}

		err := redis.Client.Do(radix.Cmd(nil, "SET", VotingEnded, "1"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		// TODO : perform the solidity transsactions to start the voting
		txhash, err := vc.StopVoting()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"link": txhash,
		})

	}
}

func GetAllVotes(vc *ethlib.VotingContractWrapper) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_ = vc
		var votes []string
		err := redis.Client.Do(radix.Cmd(&votes, "LRANGE", Votes, "0", "-1"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		// Check whether each vote is valid as per blockchain.
		// Also check whether the total number of votes is the same.
		blockchainVoteLen, err := vc.GetNumVoters()
		redisVoteLen := big.NewInt(int64(len(votes)))
		votingCompromised := false
		if err == nil && redisVoteLen.Cmp(blockchainVoteLen) != 0 {
			votingCompromised = true
		}
		for _, vote := range votes {
			if votingCompromised {
				break
			}
			verified, err := vc.VerifyVote(vote)
			// skip in case of an error... Probably, that's a valid vote.
			if err == nil && verified == false {
				votingCompromised = true
			}
		}
		if !votingCompromised {
			ctx.JSON(http.StatusOK, votes)
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Voting has been compromised",
			})

		}
	}
}
