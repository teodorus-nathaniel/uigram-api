package posts

import (
	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/auth"
)

func Routes(router *gin.RouterGroup) {
	router.GET("/posts", getPostsHandler)
	router.GET("/posts/:id", getPostHandler)
	router.GET("/posts/:id/feeds", auth.Protect(), getPostsFeedsHandler)
	router.POST("/posts", auth.Protect(), postPostHandler)
	router.GET("/users/:id/posts", getUserPostHandler)
}
