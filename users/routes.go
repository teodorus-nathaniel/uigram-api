package users

import (
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.RouterGroup) {
	router.POST("/login", loginHandler)
	router.POST("/register", registerHandler)
	router.POST("/check-user", checkUserHandler)
	router.GET("/users/:id", getUserHandler)
	router.GET("/users/:id/following", getFollowing)
	router.GET("/users/:id/followers", getFollowers)
	router.PATCH("/users/:id/add-saved", Protect(), addSavedPost)
	router.PATCH("/users/:id/delete-saved", Protect(), deleteSavedPost)
	router.PATCH("/users/:id/follow", Protect(), followUserHandler)
	router.PATCH("/users/:id/unfollow", Protect(), unfollowUserHandler)
}
