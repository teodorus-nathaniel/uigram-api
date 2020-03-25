package posts

import (
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.RouterGroup) {
	router.GET("/", getPostsHandler)
	router.GET("/:id", getPostHandler)
	router.POST("/", postPostHandler)
}
