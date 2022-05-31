package routers

import "github.com/gin-gonic/gin"

func createSignUpInRoutes(router *gin.Engine) {
	group := router.Group("/who/login")

	group.POST("/google")
	group.POST("/apple")
	group.POST("/email")
	group.POST("/phone")
	group.POST("/anon")
	group.POST("/meta")
	group.POST("/twitter")
	group.POST("/github")
	group.POST("/gitlab")
}
