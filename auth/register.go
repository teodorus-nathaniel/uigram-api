package auth

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
	"github.com/teodorus-nathaniel/uigram-api/users"
)

func registerHandler(c *gin.Context) {
	var data *users.User
	json.NewDecoder(c.Request.Body).Decode(&data)

	err := data.ValidateData()
	if err != nil {
		c.JSON(http.StatusBadRequest, jsend.GetJSendFail(err.Error()))
		return
	}

	user, err := users.InsertUser(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	token, err := getToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
	}

	user.DeriveAttributesAndHideCredentials()
	c.JSON(http.StatusCreated, jsend.GetJSendSuccess(gin.H{"user": user, "token": token}))
}