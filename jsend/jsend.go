package jsend

import (
	"github.com/gin-gonic/gin"
)

func GetJSendSuccess(data gin.H) gin.H {
	return gin.H{
		"status": "success",
		"data":   data,
	}
}

func GetJSendFail(err string) gin.H {
	return gin.H{
		"status":  "fail",
		"message": err,
	}
}
