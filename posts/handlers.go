package posts

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
	"github.com/teodorus-nathaniel/uigram-api/nodejs"
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

func getPostsHandler(c *gin.Context) {
	sort, limit, page := utils.GetQueryStringsForPagination(c)

	posts, err := getPosts(sort, limit, page, getUserFromMiddleware(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"posts": posts}))
}

func getPostsFeedsHandler(c *gin.Context) {
	_, limit, page := utils.GetQueryStringsForPagination(c)

	data, _ := c.Get("user")
	user := data.(*users.User)

	followingIds := user.Following
	posts, err := getPostsByUserId(followingIds, limit, page, getUserFromMiddleware(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"posts": posts}))
}

func getPostHandler(c *gin.Context) {
	id := c.Param("id")
	post, err := getPost(id, getUserFromMiddleware(c))

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

	res, err := insertPost(post, getUserFromMiddleware(c))

	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"post": res}))
}

func getUserPostHandler(c *gin.Context) {
	id := c.Param("id")
	sort, limit, page := utils.GetQueryStringsForPagination(c)

	posts, err := GetPostByOwner(id, sort, limit, page, getUserFromMiddleware(c))
	if err != nil {
		c.JSON(http.StatusNotFound, jsend.GetJSendFail("Fail fetching user posts"))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"posts": posts}))
}

func getUserSavedPostHandler(c *gin.Context) {
	sort, limit, page := utils.GetQueryStringsForPagination(c)

	posts, err := getSavedPosts(sort, limit, page, getUserFromMiddleware(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail("Fail fetching user posts"))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"posts": posts}))
}

type screenshotReqBody struct {
	URL string `json:"url"`
}

func postScreenshot(c *gin.Context) {
	var url screenshotReqBody
	json.NewDecoder(c.Request.Body).Decode(&url)

	matched, err := regexp.MatchString(`^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$`, url.URL)
	if err != nil || !matched {
		c.JSON(http.StatusBadRequest, jsend.GetJSendFail("Invalid URL"))
		return
	}

	if !strings.HasPrefix(url.URL, "http://") && !strings.HasPrefix(url.URL, "https://") {
		url.URL = "http://" + url.URL
	}

	res := nodejs.ExecScreenshot(url.URL)

	res = "http://localhost:8080/" + res

	c.JSON(http.StatusCreated, jsend.GetJSendSuccess(gin.H{"url": res}))
}

type patchLikesReqBody struct {
	Like    *bool `json:"like,omitempty"`
	Dislike *bool `json:"dislike,omitempty"`
}

func patchLikes(c *gin.Context) {
	user := getUserFromMiddleware(c)
	id, _ := c.Params.Get("id")
	var body patchLikesReqBody
	json.NewDecoder(c.Request.Body).Decode(&body)

	post, err := updatePostLikes(user, id, body.Like, body.Dislike)
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"post": post}))
}
