package users

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
	"github.com/teodorus-nathaniel/uigram-api/utils"
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

	res, err := AppendFollowing(user.ID.Hex(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"modifiedCount": res.ModifiedCount}))
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

	res, err := PullFollowing(user.ID.Hex(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"modifiedCount": res.ModifiedCount}))
}

func addSavedPost(c *gin.Context) {
	var object utils.ObjectId
	json.NewDecoder(c.Request.Body).Decode(&object)

	user := getUserFromMiddleware(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, jsend.GetJSendFail("You need to login first!"))
		return
	}
	res, err := AddSavedPostDatabase(user.ID.Hex(), object.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"modifiedCount": res.ModifiedCount}))
}

func deleteSavedPost(c *gin.Context) {
	var object utils.ObjectId
	json.NewDecoder(c.Request.Body).Decode(&object)

	user := getUserFromMiddleware(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, jsend.GetJSendFail("You need to login first!"))
		return
	}

	res, err := DeleteSavedPostDatabase(user.ID.Hex(), object.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"modifiedCount": res.ModifiedCount}))
}

func getFollowersOrFollowing(c *gin.Context, callback func(id string, user *User) ([]User, error)) {
	id := c.Param("id")
	user := getUserFromMiddleware(c)

	var users []User
	var err error
	if id == "self" {
		if user == nil {
			c.JSON(http.StatusUnauthorized, jsend.GetJSendFail("You need to login first!"))
			return
		}
		users, err = callback(user.ID.Hex(), user)
	} else {
		users, err = callback(id, user)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail("sorry, we couldn't get your data :("))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"users": users}))
}

func getFollowers(c *gin.Context) {
	getFollowersOrFollowing(c, getUserFollowers)
}

func getFollowing(c *gin.Context) {
	getFollowersOrFollowing(c, getUserFollowing)
}

func updateUserHandler(c *gin.Context) {
	user := getUserFromMiddleware(c)
	profilePic, _ := c.FormFile("profilePic")
	username, _ := c.GetPostForm("username")
	fullname, _ := c.GetPostForm("fullname")
	status, _ := c.GetPostForm("status")

	path := ""
	if profilePic != nil {
		path = "img/" + user.ID.Hex() + filepath.Ext(profilePic.Filename)
		err := c.SaveUploadedFile(profilePic, path)

		path = utils.URL + path

		if err != nil {
			c.JSON(http.StatusInternalServerError, jsend.GetJSendFail("File upload fail"))
			return
		}
	}

	_, err := updateUser(user.ID, username, fullname, status, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail("Update user failed"))
		return
	}

	newUser, err := GetUserById(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail("Can't get user"))
		return
	}

	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"user": newUser}))
}
