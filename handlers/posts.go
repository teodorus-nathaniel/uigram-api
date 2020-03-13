package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/repositories"
	"go.mongodb.org/mongo-driver/bson"
)

func Posts(c *gin.Context)  {
	
	posts := repositories.GetPosts(bson.M{})

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"posts": posts,
		},
	})
}