package routers

import (
	"github.com/gin-gonic/gin"
)

func createAuthenticationRoutes(router *gin.Engine) {
	createLoginInRoutes(router)
	createSignUpInRoutes(router)
	createGeneralUserRoutes(router)
}

func createAdminRoutes(router *gin.Engine) {
	createAdminBaseRoutes(router)
	createAdminSettingsRoutes(router)
}

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
