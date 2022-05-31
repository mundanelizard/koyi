package routers

import "github.com/gin-gonic/gin"

func createMFARoutes(router *gin.Engine) {
	group := router.Group("/mfa")

	group.POST("/totp") // email and phone number
	group.POST("/totp")
	group.GET("/otp/verify")
	group.GET("/totp/verify")
}
