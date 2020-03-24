package posts

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/helper"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
	"go.mongodb.org/mongo-driver/bson"
)

func GetPostsHandler(c *gin.Context) {
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

	posts, err := GetPosts(sort, limit, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess("posts", posts))
}

func GetPostHandler(c *gin.Context) {
	id := c.Param("id")
	post, err := GetPost(bson.M{
		"_id": id,
	})

	if err != nil || post == nil {
		c.JSON(http.StatusNotFound, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess("post", post))
}

func PostPostHandler(c *gin.Context) {
	var post Post
	json.NewDecoder(c.Request.Body).Decode(&post)

	bsonData, _ := helper.ConvertStructToBSON(post)

	res, err := InsertPost(bsonData)

	if err != nil || res == nil {
		c.JSON(http.StatusBadRequest, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess("post", res))
}
