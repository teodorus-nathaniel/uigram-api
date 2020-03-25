package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

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

	routerGroup := router.Group("/api/v1")
	initializeRoutes(routerGroup)

	fmt.Println("Server started...")
	router.Run()
}
