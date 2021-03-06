package users

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
)

type checkUserBody struct {
	Token *string `json:"token"`
}

func checkUserHandler(c *gin.Context) {
	var token *checkUserBody
	json.NewDecoder(c.Request.Body).Decode(&token)

	user, err := ValidateToken(token.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, jsend.GetJSendFail(err.Error()))
		return
	}

	newToken, err := getToken(user.ID.Hex())

	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
		return
	}

	user.DeriveAttributesAndHideCredentials(getUserFromMiddleware(c))
	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"user": user, "token": newToken}))
}
