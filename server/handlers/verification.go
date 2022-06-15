package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mundanelizard/koyi/server/handlers/middlewares"
	"log"
	"net/http"
)

func verifyHandler(c *gin.Context) {
	var intentId string
	var code string
	var err error

	if c.Request.Method == http.MethodPost {
		var body map[string]string
		err = c.BindJSON(body)
	} else {
		intentId = c.Param("intentId")
		code = c.Param("code")
	}

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
	}

	ok, err := verifyUser(intentId, code)

	if !ok {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{})
}

func CreateVerificationRoutes(router *gin.RouterGroup) {
	group := router.Group("/auth/verify")

	group.POST("/", middlewares.ValidateEmailSignUp, emailSignUpHandler)
	group.GET("/:intentId/:code", middlewares.ValidatePhoneNumberSignUp, phoneNumberSignUpHandler)
}
