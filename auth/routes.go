package auth

import "github.com/gin-gonic/gin"

func Routes(router *gin.RouterGroup) {
	router.POST("/login", loginHandler)
	router.POST("/register", registerHandler)
	router.POST("/check-user", checkUserHandler)
}
