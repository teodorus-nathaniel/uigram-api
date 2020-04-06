package posts

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
	"github.com/teodorus-nathaniel/uigram-api/users"
)

func getQueryStringsForPagination(c *gin.Context) (string, int, int) {
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

	return sort, limit, page
}

func getPostsHandler(c *gin.Context) {
	sort, limit, page := getQueryStringsForPagination(c)

	posts, err := getPosts(sort, limit, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"posts": posts}))
}

func getPostsFeedsHandler(c *gin.Context) {
	_, limit, page := getQueryStringsForPagination(c)

	data, _ := c.Get("user")
	user := data.(*users.User)

	followingIds := user.Following
	posts, err := getPostsByUserId(followingIds, limit, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"posts": posts}))
}

func getPostHandler(c *gin.Context) {
	id := c.Param("id")
	post, err := getPost(id)

	if err != nil {
		c.JSON(http.StatusNotFound, jsend.GetJSendFail("post with id: "+id+" was not found"))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"post": post}))
}

func postPostHandler(c *gin.Context) {
	var post Post
	json.NewDecoder(c.Request.Body).Decode(&post)

	err := post.validateData()
	if err != nil {
		c.JSON(http.StatusBadRequest, jsend.GetJSendFail(err.Error()))
		return
	}

	post.fillEmptyValues()

	user, _ := c.Get("user")
	post.UserID = user.(*users.User).ID.Hex()

	res, err := insertPost(post)

	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"post": res}))
}
