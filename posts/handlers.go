package posts

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getPostsHandler(c *gin.Context) {
	sort := c.Query("sort")
	limitTemp := c.Query("limit")
	pageTemp := c.Query("page")

	limit, _ := strconv.Atoi(limitTemp)
	page, _ := strconv.Atoi(pageTemp)

	if sort == "" {
		sort = "timestamp"
	}
	if limit == 0 {
		limit = 15
	}
	if page == 0 {
		page = 1
	}

	posts, err := getPosts(sort, limit, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess("posts", posts))
}

func getPostHandler(c *gin.Context) {
	id := c.Param("id")
	post, err := getPost(id)

	if err != nil {
		c.JSON(http.StatusNotFound, jsend.GetJSendFail("post with id: "+id+" was not found"))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess("post", post))
}

func postPostHandler(c *gin.Context) {
	var post Post
	json.NewDecoder(c.Request.Body).Decode(&post)

	post.ID = primitive.NilObjectID
	err := post.validateData()
	if err != nil {
		c.JSON(http.StatusBadRequest, jsend.GetJSendFail(err.Error()))
		return
	}

	res, err := insertPost(post)

	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess("post", res))
}
