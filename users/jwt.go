package users

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
)

var jwtKey []byte

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading env", err.Error())
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
	expireTime := 7 * 24 * time.Hour
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

func ValidateToken(token *string) (*User, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(*token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	claims, ok := tkn.Claims.(*Claims)
	if !ok || !tkn.Valid {
		return nil, errors.New("Token invalid or has expired. Please login first")
	}

	user, err := GetUserById(claims.ID)
	if err != nil {
		return nil, errors.New("user for this token is invalid")
	}

	return user, err
}

func Protect() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("user")
		if user == nil {
			c.JSON(http.StatusUnauthorized, jsend.GetJSendFail("The resource you are looking for is restricted. Please login first"))
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, exists := getTokenFromHeader(c)

		if !exists {
			c.Set("user", nil)
			return
		}

		user, err := ValidateToken(token)
		if err != nil {
			c.Set("user", nil)
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
