package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/teodorus-nathaniel/uigram-api/database"
	"github.com/teodorus-nathaniel/uigram-api/posts"
)

func initializeRoutes(router *gin.RouterGroup) {
	posts.Routes(router.Group("/posts"))
}

func main() {
	defer database.Client.Disconnect(database.Context)
	router := gin.Default()

	routerGroup := router.Group("/api/v1")
	initializeRoutes(routerGroup)

	fmt.Println("Server started...")
	router.Run()
}
