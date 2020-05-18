package comments

import (
	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/users"
)

func Routes(router *gin.RouterGroup) {
	router.POST("/comments", users.Protect(), postComment)
	router.POST("/comments/:id/replies", users.Protect(), postReply)
	router.GET("/posts/:id/comments", getCommentsByPostIdHandler)
	router.GET("/posts/:id/comments/:commentId/replies", getCommentsReplies)
	router.PATCH("/posts/:id/comments/:commentId/likes", users.Protect(), patchCommentLikes)
}
