package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/ui-share-api/data"
	_ "github.com/teodorus-nathaniel/ui-share-api/database"
)

func GetAllPosts(c *gin.Context) {
	// sort := c.Query("sort")
	// database.Client.Database()
	c.JSON(http.StatusOK, data.DummyPosts)
}
