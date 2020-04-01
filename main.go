package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"

	"github.com/teodorus-nathaniel/uigram-api/auth"
	"github.com/teodorus-nathaniel/uigram-api/database"
	"github.com/teodorus-nathaniel/uigram-api/posts"
	"github.com/teodorus-nathaniel/uigram-api/users"
)

func initializeRoutes(router *gin.RouterGroup) {
	auth.Routes(router.Group("/"))
	posts.Routes(router.Group("/posts"))
	users.Routes(router.Group("/users"))
}

func main() {
	defer database.Client.Disconnect(database.Context)
	router := gin.Default()

	router.Use(cors.New(cors.Options{
		Debug:          true,
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}))

	routerGroup := router.Group("/api/v1")
	initializeRoutes(routerGroup)

	fmt.Println("Server started...")
	router.Run()
}
