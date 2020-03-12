package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/ui-share-api/handlers"
)

func UsePostsRoutes(router *gin.RouterGroup) {
	postsGroup := router.Group("/posts")
	postsGroup.GET("/", handlers.GetAllPosts)
}
