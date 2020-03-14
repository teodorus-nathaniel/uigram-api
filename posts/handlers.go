package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
	"go.mongodb.org/mongo-driver/bson"
)

func GetPostsHandler(c *gin.Context) {
	posts := GetPosts(bson.M{})

	c.JSON(http.StatusOK, jsend.GetJSendSuccess("posts", posts))
}


func GetPostHandler(c *gin.Context){
	id := c.Param("id")
	post, err := GetPost(bson.M{
		"_id" : id,
	})

	if err != nil || post == nil{
		c.JSON(http.StatusNotFound,jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess("post", post))
}	