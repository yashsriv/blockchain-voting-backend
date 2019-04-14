package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jlaffaye/ftp"
	"github.com/spf13/viper"
)

type loginRequest struct {
	Username string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type jwtClaims struct {
	Username    string `json:"username"`
	IsAdmin     bool   `json:"isAdmin"`
	IsCandidate bool   `json:"isCandidate"`
	jwt.StandardClaims
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
	ftpEnabled := viper.GetBool("ftp.enabled")

	secret := []byte(viper.GetString("jwt.secret"))

	admin := viper.GetString("admin.username")
	candidatesList := viper.GetStringSlice("admin.candidates")

	var candidates = make(map[string]bool)
	for _, candidate := range candidatesList {
		candidates[candidate] = true
	}

	return func(ctx *gin.Context) {
		var request loginRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if ftpEnabled {
			valid, err := request.verify(ftpHost, ftpPort)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if !valid {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}

		username := request.Username
		currentTime := time.Now()

		_, isCandidate := candidates[username]
		claims := jwtClaims{
			username,
			username == admin,
			isCandidate,
			jwt.StandardClaims{
				IssuedAt:  currentTime.Unix(),
				ExpiresAt: currentTime.Add(time.Hour).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(secret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"token":       tokenString,
			"username":    username,
			"isAdmin":     username == admin,
			"isCandidate": isCandidate,
		})
	}
}

func AuthCheck() gin.HandlerFunc {
	secret := []byte(viper.GetString("jwt.secret"))
	return func(ctx *gin.Context) {
		splitted := strings.Split(ctx.GetHeader("Authorization"), "Bearer ")
		if len(splitted) == 1 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "credentials missing"})
			return
		}

		reqToken := splitted[1]

		token, err := jwt.ParseWithClaims(reqToken, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":  "credentials rejected",
				"reason": err.Error(),
			})
			return
		}

		claims, ok := token.Claims.(*jwtClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":  "credentials rejected",
				"reason": "claims object doesn't match expected",
			})
			return
		}

		ctx.Set("username", claims.Username)
		ctx.Set("isAdmin", claims.IsAdmin)
		ctx.Set("isCandidate", claims.IsCandidate)

		ctx.Next()
	}
}
