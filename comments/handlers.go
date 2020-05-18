package comments

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
	"github.com/teodorus-nathaniel/uigram-api/users"
	"github.com/teodorus-nathaniel/uigram-api/utils"
)

func getUserFromMiddleware(c *gin.Context) *users.User {
	user, _ := c.Get("user")
	if user == nil {
		return nil
	}
	return user.(*users.User)
}

func getCommentsByPostIdHandler(c *gin.Context) {
	postID, _ := c.Params.Get("id")
	sort, limit, page := utils.GetQueryStringsForPagination(c)

	comments, err := getCommentsByPostId(postID, sort, limit, page, getUserFromMiddleware(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"comments": comments}))
}

func getCommentsReplies(c *gin.Context) {
	commentID, _ := c.Params.Get("commentId")
	sort, limit, page := utils.GetQueryStringsForPagination(c)

	replies, err := getCommentsRepliesByParent(commentID, sort, limit, page, getUserFromMiddleware(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"replies": replies}))
}

type patchLikesReqBody struct {
	Like    *bool `json:"like,omitempty"`
	Dislike *bool `json:"dislike,omitempty"`
}

func patchCommentLikes(c *gin.Context) {
	user := getUserFromMiddleware(c)
	id, _ := c.Params.Get("commentId")

	var body patchLikesReqBody
	json.NewDecoder(c.Request.Body).Decode(&body)

	updatedComment, err := updateCommentLikes(id, body.Like, body.Dislike, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"comment": updatedComment}))
}

type PostCommentBody struct {
	PostID  string `json:"postId"`
	Content string `json:"content"`
}

func postComment(c *gin.Context) {
	var body PostCommentBody
	json.NewDecoder(c.Request.Body).Decode(&body)

	if body.Content == "" {
		c.JSON(http.StatusBadRequest, jsend.GetJSendFail("Comment can't be empty"))
		return
	}

	comment, err := insertComment(body.PostID, body.Content, getUserFromMiddleware(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, jsend.GetJSendSuccess(gin.H{"comment": comment}))
}

func postReply(c *gin.Context) {
	var body PostCommentBody
	json.NewDecoder(c.Request.Body).Decode(&body)
	id, _ := c.Params.Get("id")

	if body.Content == "" {
		c.JSON(http.StatusBadRequest, jsend.GetJSendFail("Comment can't be empty"))
		return
	}

	reply, err := insertReply(body.PostID, id, body.Content, getUserFromMiddleware(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, jsend.GetJSendSuccess(gin.H{"reply": reply}))

}
