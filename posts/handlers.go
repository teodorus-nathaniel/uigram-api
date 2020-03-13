package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func PostsHandler(c *gin.Context) {
	posts := GetPosts(bson.M{})

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"posts": posts,
		},
	})
}
