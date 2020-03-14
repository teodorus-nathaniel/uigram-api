package jsend

import (
	"github.com/gin-gonic/gin"
)

func GetJSendSuccess(name string, data interface{}) gin.H {
	return gin.H{
		"status": "success",
		"data": gin.H{
			name: data,
		},
	}
}

func GetJSendFail(err string) gin.H {
	return gin.H{
		"status": "fail",
		"message": err,
	}
}
