package handlers

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jlaffaye/ftp"
	"github.com/spf13/viper"
	"golang.org/x/crypto/sha3"
)

type loginRequest struct {
	Username string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func (req *loginRequest) verify(host string, port string) (bool, error) {
	conn, err := ftp.Dial(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Printf("[warn] %v\n", err)
		return false, err
	}
	defer conn.Quit()

	err = conn.Login(req.Username, req.Password)
	if err != nil {
		errStr := err.Error()
		codeStr := strings.Split(errStr, " ")[0]
		if code, errConv := strconv.Atoi(codeStr); errConv == nil && code == 530 {
			return false, nil
		} else {
			log.Printf("[warn] %v\n", err)
		}
	}
	defer conn.Logout()

	return true, nil
}

func HandleLogin() gin.HandlerFunc {
	ftpHost := viper.GetString("ftp.host")
	ftpPort := viper.GetString("ftp.port")
	secret := viper.GetString("cookie.secret")

	return func(ctx *gin.Context) {
		var request loginRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		valid, err := request.verify(ftpHost, ftpPort)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		username := request.Username
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		hashValue := []byte(fmt.Sprintf("%s:%s:%s", username, timestamp, secret))
		hasher := sha3.New256()
		hasher.Write(hashValue)
		auth := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

		ctx.SetCookie("username", username, 0, "", "", false, false)
		ctx.SetCookie("timestamp", timestamp, 0, "", "", false, false)
		ctx.SetCookie("auth", auth, 0, "", "", false, false)
	}
}

func AuthCheck() gin.HandlerFunc {
	secret := viper.GetString("cookie.secret")
	return func(ctx *gin.Context) {
		username, err := ctx.Cookie("username")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		timestamp, err := ctx.Cookie("timestamp")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		auth, err := ctx.Cookie("auth")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashValue := []byte(fmt.Sprintf("%s:%s:%s", username, timestamp, secret))
		hasher := sha3.New256()
		hasher.Write(hashValue)
		expectedAuth := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

		if auth != expectedAuth {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "credentials rejected"})
			return
		}

		ctx.Next()
	}
}
