package posts

import (
	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/auth"
)

func Routes(router *gin.RouterGroup) {
	router.GET("/", getPostsHandler)
	router.GET("/:id", getPostHandler)
	router.GET("/:id/feeds", auth.Protect(), getPostsFeedsHandler)
	router.POST("/", auth.Protect(), postPostHandler)
}
