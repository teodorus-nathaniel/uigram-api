package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"

	"github.com/teodorus-nathaniel/uigram-api/database"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
	"github.com/teodorus-nathaniel/uigram-api/posts"
	"github.com/teodorus-nathaniel/uigram-api/users"
)

func initializeRoutes(router *gin.RouterGroup) {
	posts.Routes(router)
	users.Routes(router)
}

func main() {
	defer database.Client.Disconnect(database.Context)
	router := gin.Default()

	router.GET("/img/:filename", func(c *gin.Context) {
		filePath := "img/"
		file, exists := c.Params.Get("filename")
		if !exists || !strings.HasSuffix(file, ".jpg") {
			c.JSON(http.StatusBadRequest, jsend.GetJSendFail("file not found"))
		}
		c.File(filePath + file)
	})

	router.Use(cors.New(cors.Options{
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
