package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/ui-share-api/data"
	"github.com/teodorus-nathaniel/ui-share-api/database"
)

func GetAllPosts(c *gin.Context) {
	// sort := c.Query("sort")
	database.Database.Collection("posts").InsertOne(database.Context, data.DummyPosts[0])
	c.JSON(http.StatusOK, data.DummyPosts)
}
