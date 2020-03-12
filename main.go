package main

import (
	"fmt"
	_ "fmt"

	_ "net/http"

	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/ui-share-api/routes"
)

func initializeRoutes(router *gin.RouterGroup) {
	routes.UsePostsRoutes(router)
}

func main() {
	router := gin.Default()

	routerGroup := router.Group("/api/v1")
	initializeRoutes(routerGroup)

	fmt.Println("Server started...")

	router.Run()
}
