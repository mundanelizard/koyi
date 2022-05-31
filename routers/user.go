package routers

import "github.com/gin-gonic/gin"

func createGeneralUserRoutes(router *gin.Engine) {
	group := router.Group("/who")

	group.POST("/logout")
	group.POST("/reset") // { type: "", id: "" }

	group.PATCH("/metadata")
	group.PUT("/password") // { type: "", newPassword: "", oldPassword: "" }

	group.GET("/")
	group.DELETE("/")
}
