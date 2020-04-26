package users

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
)

func loginHandler(c *gin.Context) {
	credentials := &Credentials{}
	json.NewDecoder(c.Request.Body).Decode(&credentials)

	err := credentials.ValidateEmailPassword()
	if err != nil {
		c.JSON(http.StatusBadRequest, jsend.GetJSendFail(err.Error()))
		return
	}

	user, err := GetUserByEmail(credentials.Email)

	isRightPass := false
	if user != nil {
		isRightPass = user.IsRightPassword(credentials.Password)
	}
	if err != nil || !isRightPass {
		c.JSON(http.StatusUnauthorized, jsend.GetJSendFail("invalid email or password"))
		return
	}

	token, err := getToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsend.GetJSendFail(err.Error()))
	}

	user.DeriveAttributesAndHideCredentials(getUserFromMiddleware(c))
	c.JSON(http.StatusOK, jsend.GetJSendSuccess(gin.H{"user": user, "token": token}))
}
