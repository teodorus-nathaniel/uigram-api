package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/teodorus-nathaniel/ui-share-api/database"
	"github.com/teodorus-nathaniel/ui-share-api/routes"
)

func initializeRoutes(router *gin.RouterGroup) {
	routes.UsePostsRoutes(router)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading env", err.Error())
	}

	database.InitializeClient()

	router := gin.Default()

	routerGroup := router.Group("/api/v1")
	initializeRoutes(routerGroup)

	fmt.Println("Server started...")

	router.Run()
}
