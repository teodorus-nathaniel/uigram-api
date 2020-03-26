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

func sendTokenAsCookie(c *gin.Context, id string) {
	expireTime := 1 * time.Minute
	token, err := createToken(id, expireTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}
	c.SetCookie("token", *token, int(expireTime), "/", "localhost:8080", false, false)
}

func Protect() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")

		if err != nil {
			c.JSON(http.StatusUnauthorized, jsend.GetJSendFail("The resource you are looking for is restricted. Please login first"))
		}

		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		claims, ok := tkn.Claims.(*Claims)
		if !ok || !tkn.Valid {
			c.JSON(http.StatusUnauthorized, jsend.GetJSendFail("Token invalid or has expired. Please login first"))
		}

		user, err := users.GetUserById(claims.ID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, jsend.GetJSendFail("user for this token is invalid"))
		}

		c.Set("user", user)
		c.Next()
	}
}
