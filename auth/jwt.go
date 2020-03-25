package auth

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
)

var jwtKey []byte

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading env", err.Error())
	}

	key := os.Getenv("JWT_SECRET")
	jwtKey = []byte(key)
}

func createToken(id string, expireTime time.Duration) (*string, error) {
	expirationTime := time.Now().Add(expireTime)
	claims := &Claims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func sendTokenAsCookie(c *gin.Context, id string) {
	expireTime := 15 * time.Minute
	token, err := createToken(id, expireTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}
	c.SetCookie("token", *token, int(expireTime), "/", "localhost:8080", false, false)
}
