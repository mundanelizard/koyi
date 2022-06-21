package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mundanelizard/koyi/server/handlers"
	"log"
)

func createVersion1Handlers(engine *gin.Engine) *gin.Engine {
	router := engine.Group("v1")
	handlers.CreateSignUpRoutes(router)       // give users roles on signup
	handlers.CreateVerificationRoutes(router) // use this for otp to
	handlers.CreateSignInRoutes(router)
	// handlers.CreatePasswordRoutes(router)
	// handlers.CreateMetadataRoutes(router)
	// handlers.CreateUserRoutes(router)
	return engine
}

// SetUpServer sets up the server for handling all post request.
func setUpServer() *gin.Engine {
	/**
	- Keep user action logs.
	*/
	engine := gin.Default()

	return createVersion1Handlers(engine)
}

func main() {
	server := setUpServer()
	err := server.Run()
	if err != nil {
		log.Fatalf("Finally, there's an error %s", err)
	}
}

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
*/
