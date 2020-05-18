package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"

	"github.com/teodorus-nathaniel/uigram-api/comments"
	"github.com/teodorus-nathaniel/uigram-api/database"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
	"github.com/teodorus-nathaniel/uigram-api/posts"
	"github.com/teodorus-nathaniel/uigram-api/users"
)

func initializeRoutes(router *gin.RouterGroup) {
	posts.Routes(router)
	users.Routes(router)
	comments.Routes(router)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	defer database.Client.Disconnect(database.Context)
	router := gin.Default()

	router.GET("/img/:filename", func(c *gin.Context) {
		filePath := "img/"
		file, exists := c.Params.Get("filename")
		if !exists || (!strings.HasSuffix(file, ".jpg") && (!strings.HasSuffix(file, ".jpeg")) && !strings.HasSuffix(file, ".png")) {
			c.JSON(http.StatusBadRequest, jsend.GetJSendFail("file not found"))
			return
		}
		c.File(filePath + file)
	})

	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://uigram.herokuapp.com", "https://uigram.herokuapp.com"},
		Debug:          true,
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		AllowedMethods: []string{"POST", "GET", "HEAD", "PATCH"},
	}))

	routerGroup := router.Group("/api/v1")
	routerGroup.Use(users.GetUserMiddleware())
	initializeRoutes(routerGroup)

	fmt.Println("Server started...")
	router.Run()
}
