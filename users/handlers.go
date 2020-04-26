package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
)

func getUserFromMiddleware(c *gin.Context) *User {
	user, _ := c.Get("user")
	if user == nil {
		return nil
	}
	return user.(*User)
}

func getUserHandler(c *gin.Context) {
	id := c.Param("id")
	user, err := GetUserById(id)

	if err != nil {
		c.JSON(http.StatusNotFound, jsend.GetJSendFail("User with id: "+id+" was not found"))
		return
	}

	user.DeriveAttributesAndHideCredentials(getUserFromMiddleware(c))

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"user": user}))
}

func followUserHandler(c *gin.Context) {
	id := c.Param("id")
	data, _ := c.Get("user")

	user := data.(*User)
	_, err := AppendFollower(id, user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	_, err = AppendFollowing(user.ID.Hex(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func unfollowUserHandler(c *gin.Context) {
	id := c.Param("id")
	data, _ := c.Get("user")

	user := data.(*User)
	_, err := PullFollower(id, user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	_, err = PullFollowing(user.ID.Hex(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
