package users

import (
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.RouterGroup) {
	router.POST("/login", loginHandler)
	router.POST("/register", registerHandler)
	router.POST("/check-user", checkUserHandler)
	router.GET("/users/:id", getUserHandler)
	router.PATCH("/users/:id/follow", Protect(), followUserHandler)
	router.PATCH("/users/:id/unfollow", Protect(), unfollowUserHandler)
}
