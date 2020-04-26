package posts

import (
	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/users"
)

func Routes(router *gin.RouterGroup) {
	router.GET("/posts", getPostsHandler)
	router.GET("/posts/:id", getPostHandler)
	router.GET("/users/:id/posts", getUserPostHandler)
	router.GET("/posts/:id/feeds", users.Protect(), getPostsFeedsHandler)
	router.POST("/posts", users.Protect(), postPostHandler)
}
