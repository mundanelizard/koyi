package routers

import "github.com/gin-gonic/gin"

func createAdminSettingsRoutes(router *gin.Engine) {
	group := router.Group("/who/admins")

	group.GET("/settings/providers")
	group.POST("/settings/providers")

	// sms, verification email, password rest, email address change
	group.GET("/settings/templates")
	group.POST("/settings/templates")

	// change your smtp provider (AWS SES or Any SMTP)
	group.GET("/settings/stmp")
	group.POST("/settings/stmp")

	group.GET("/settings/sms")
	group.POST("/settings/sms")
}

