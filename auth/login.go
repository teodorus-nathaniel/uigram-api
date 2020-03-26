package auth

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teodorus-nathaniel/uigram-api/jsend"
	"github.com/teodorus-nathaniel/uigram-api/users"
)

func loginHandler(c *gin.Context) {
	credentials := &users.Credentials{}
	json.NewDecoder(c.Request.Body).Decode(&credentials)
	// credentials.Email = c.Query("email")
	// credentials.Password = c.Query("password")

	err := credentials.ValidateEmailPassword()
	if err != nil {
		c.JSON(http.StatusBadRequest, jsend.GetJSendFail(err.Error()))
		return
	}

	user, err := users.GetUserByEmailAndPassword(credentials.Email, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, jsend.GetJSendFail("invalid email or password"))
		return
	}

	sendTokenAsCookie(c, user.ID.Hex())

	c.JSON(http.StatusOK, jsend.GetJSendSuccess("user", user))
}
