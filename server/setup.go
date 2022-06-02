package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Groupable is an interface for *gin.Group
type Groupable interface {
	Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup
}

func createAdminBaseRoutes(router Groupable) {
	group := router.Group("/admins", noOperationHandler)

	group.GET("/dashboard", noOperationHandler)

	group.GET("/", noOperationHandler)         // gets all members
	group.POST("/", noOperationHandler)        // adds new members
	group.PUT("/:adminId", noOperationHandler) // gets the current admin details

	group.GET("/users", noOperationHandler) // filter?roles, last-logged-in,
	group.GET("/users/:userId", noOperationHandler)
}

func createAdminSettingsRoutes(router Groupable) {
	group := router.Group("/who/admins", noOperationHandler)

	group.GET("/settings/providers", noOperationHandler)
	group.POST("/settings/providers", noOperationHandler)

	// sms, verification email, password rest, email address change
	group.GET("/settings/templates", noOperationHandler)
	group.POST("/settings/templates", noOperationHandler)

	// change your smtp provider (AWS SES or Any SMTP)
	group.GET("/settings/stmp", noOperationHandler)
	group.POST("/settings/stmp", noOperationHandler)

	group.GET("/settings/sms", noOperationHandler)
	group.POST("/settings/sms", noOperationHandler)
}

func createLoginInRoutes(router Groupable) {
	group := router.Group("/who/login")

	group.POST("/google", noOperationHandler)
	group.POST("/apple", noOperationHandler)
	group.POST("/email", noOperationHandler)
	group.POST("/email/totp", noOperationHandler)
	group.POST("/phone", noOperationHandler)
	group.POST("/phone/totp", noOperationHandler)
	group.POST("/anon", noOperationHandler)
	group.POST("/meta", noOperationHandler)
	group.POST("/twitter", noOperationHandler)
	group.POST("/github", noOperationHandler)
	group.POST("/gitlab", noOperationHandler)

}

func createMFARoutes(router *gin.Engine) {
	group := router.Group("/mfa")

	group.POST("/totp", noOperationHandler) // email and phone number
	group.POST("/totp", noOperationHandler)
	group.GET("/otp/verify", noOperationHandler)
	group.GET("/totp/verify", noOperationHandler)
}

func createSignUpInRoutes(router Groupable) {
	group := router.Group("/who/login")

	group.POST("/google", noOperationHandler)
	group.POST("/apple", noOperationHandler)
	group.POST("/email", noOperationHandler)
	group.POST("/phone", noOperationHandler)
	group.POST("/anon", noOperationHandler)
	group.POST("/meta", noOperationHandler)
	group.POST("/twitter", noOperationHandler)
	group.POST("/github", noOperationHandler)
	group.POST("/gitlab", noOperationHandler)
}

func createGeneralUserRoutes(router Groupable) {
	group := router.Group("/who")

	group.POST("/logout", noOperationHandler)
	group.POST("/reset", noOperationHandler) // { type: "", id: "" }

	group.PATCH("/metadata", noOperationHandler)
	group.PUT("/password", noOperationHandler) // { type: "", newPassword: "", oldPassword: "" }

	group.GET("/", noOperationHandler)
	group.DELETE("/", noOperationHandler)
}

func createAuthorizationRoutes(router *gin.Engine) {
	group := router.Group("/what", noOperationHandler)
	group.GET("/roles", noOperationHandler)
	group.GET("/", noOperationHandler)
}

// createAuthenticationRoutes handles all the authentication request made in the application.
func createAuthenticationRoutes(router Groupable) {
	createLoginInRoutes(router)
	createSignUpInRoutes(router)
	createGeneralUserRoutes(router)
}

// createAdminRoutes handles the application administration requests.
func createAdminRoutes(router Groupable) {
	createAdminBaseRoutes(router)
	createAdminSettingsRoutes(router)
}

// noOperationHandler handle the absence of a route implementation.
func noOperationHandler(context *gin.Context) {
	fmt.Println(context)
	context.AbortWithStatusJSON(http.StatusOK, map[string]string{"Ping": "Pong"})
}

// SetUpServer sets up the server for handling all post request.
func SetUpServer() *gin.Engine {
	/**
	- Keep user action logs.
	*/
	router := gin.Default()

	createAuthenticationRoutes(router)
	createAuthorizationRoutes(router)
	createAdminRoutes(router)
	createMFARoutes(router)

	return router
}
