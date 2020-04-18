package users

import (
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.RouterGroup) {
	router.GET("/users/:id", getUserHandler)
	// router.GET("/", GetPostsHandler)
	// router.POST("/", PostPostHandler)
}
