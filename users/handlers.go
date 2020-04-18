package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
)

func getUserHandler(c *gin.Context) {
	id := c.Param("id")
	user, err := GetUserById(id)

	if err != nil {
		c.JSON(http.StatusNotFound, jsend.GetJSendFail("User with id: "+id+" was not found"))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"user": user}))
}
