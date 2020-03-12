package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/ui-share-api/data"
)

func GetAllPosts(c *gin.Context) {
	// sort := c.Query("sort")
	c.JSON(http.StatusOK, data.DummyPosts)
}
