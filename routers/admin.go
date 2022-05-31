package routers

import "github.com/gin-gonic/gin"

func createAdminBaseRoutes(router *gin.Engine) {
	group := router.Group("/admins")

	group.GET("/dashboard")

	group.GET("/")         // gets all members
	group.POST("/")        // adds new members
	group.PUT("/:adminId") // gets the current admin details

	group.GET("/users") // filter?roles, last-logged-in,
	group.GET("/users/:userId")
}


