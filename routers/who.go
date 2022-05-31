package routers

import "github.com/gin-gonic/gin"

func createAuthorizationRoutes(router *gin.Engine) {
	group := router.Group("/what")
	group.GET("/roles")
	group.GET("/")
}
