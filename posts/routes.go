package posts

import (
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.RouterGroup) {
	router.GET("/", GetPostsHandler)
	router.GET("/:id", GetPostHandler)
	router.POST("/", PostPostHandler)
}
