package auth

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
	"github.com/teodorus-nathaniel/uigram-api/users"
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

func getToken(id string) (*string, error) {
	expireTime := 1 * time.Minute
	token, err := createToken(id, expireTime)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func getTokenFromHeader(c *gin.Context) (*string, bool) {
	header := c.GetHeader("Authorization")
	stringTokens := strings.Split(header, " ")

	if len(stringTokens) != 2 {
		return nil, false
	}
	return &stringTokens[1], true
}

func Protect() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, exists := getTokenFromHeader(c)

		if !exists {
			c.JSON(http.StatusUnauthorized, jsend.GetJSendFail("The resource you are looking for is restricted. Please login first"))
			c.Abort()
			return
		}

		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(*token, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		claims, ok := tkn.Claims.(*Claims)
		if !ok || !tkn.Valid {
			c.JSON(http.StatusUnauthorized, jsend.GetJSendFail("Token invalid or has expired. Please login first"))
			c.Abort()
			return
		}

		user, err := users.GetUserById(claims.ID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, jsend.GetJSendFail("user for this token is invalid"))
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
