package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mundanelizard/koyi/server/handlers/middlewares"
)

func verifyPhoneNumberHandler(c *gin.Context) {

}

func CreateVerificationRoutes(router *gin.RouterGroup) {
	group := router.Group("/auth/verify")

	group.POST("/", middlewares.ValidateEmailSignUp, emailSignUpHandler)
	group.GET("/:intentId/:code", middlewares.ValidatePhoneNumberSignUp, phoneNumberSignUpHandler)
}
