package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mundanelizard/koyi/server/handlers"
	"log"
)

func createVersion1Handlers(engine handlers.GroupableRoutes) {
	router := engine.Group("v1")
	handlers.CreateAuthenticationRoutes(router)
	//createAuthorizationRoutes(router)
	//createAdminRoutes(router)
	//createMFARoutes(router)
}

// SetUpServer sets up the server for handling all post request.
func setUpServer() *gin.Engine {
	/**
	- Keep user action logs.
	*/
	engine := gin.Default()

	createVersion1Handlers(engine)

	return engine
}

func main() {
	server := setUpServer()
	err := server.Run()
	if err != nil {
		log.Fatalf("Finally, there's an error %s", err)
	}
}

/*
func createGeneralUserRoutes(router Groupable) {
	group := router.Group("/who")

	group.POST("/logout", noOperationHandler)
	group.POST("/reset", noOperationHandler) // { type: "", id: "" }

	group.PATCH("/metadata", noOperationHandler)
	group.PUT("/password", noOperationHandler) // { type: "", newPassword: "", oldPassword: "" }

	group.GET("/", noOperationHandler)
	group.DELETE("/", noOperationHandler)
}

func createAuthorizationRoutes(router Groupable) {
	group := router.Group("/what", noOperationHandler)
	group.GET("/roles", noOperationHandler)
	group.GET("/", noOperationHandler)
}
*/

/*
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

	group.POST("/phone", middlewares.ValidatePhoneNumberSignIn, handlers.PhoneNumberSignInHandler)
	group.POST("/email", middlewares.ValidateEmailSignIn, handlers.EmailSignInHandler)

	group.POST("/email/totp", noOperationHandler) // logs the user in using email and otp
	group.POST("/phone/totp", noOperationHandler) // logs the user in using phone number and otp
}

func createMFARoutes(router Groupable) {
	group := router.Group("/mfa")

	group.POST("/totp", noOperationHandler) // email and phone number
	group.GET("/otp/verify", noOperationHandler)
}
*/

/*
// createAdminRoutes handles the application administration requests.
func createAdminRoutes(router Groupable) {
	// Add authorized domains feature
	createAdminBaseRoutes(router)
	createAdminSettingsRoutes(router)
}

// noOperationHandler handle the absence of a route implementation.
func noOperationHandler(context *gin.Context) {
	fmt.Println(context)
	context.AbortWithStatusJSON(http.StatusOK, map[string]string{"Ping": "Pong"})
}
*/
