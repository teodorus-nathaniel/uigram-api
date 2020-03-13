package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/teodorus-nathaniel/uigram-api/database"
	"github.com/teodorus-nathaniel/uigram-api/routes"
)

func initializeRoutes(router *gin.RouterGroup) {
	routes.UsePostsRoutes(router)
}

func main() {
	defer database.Client.Disconnect(database.Context)
	router := gin.Default()

	routerGroup := router.Group("/api/v1")
	initializeRoutes(routerGroup)

	fmt.Println("Server started...")
	router.Run()
}
