package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/mundanelizard/koyi/server/helpers"
	"log"
	"net/http"
)

func ValidateEmailSignUp(c *gin.Context) {
	var details map[string]interface{}

	if err := c.BindJSON(&details); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	email, ok := details["email"].(string)

	if ok && helpers.IsValidEmail(email) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
	}

	password, ok := details["password"].(string)

	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
	}

	err := helpers.IsValidPassword(password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
	}

	c.Set("email", email)
	c.Set("password", password)
	c.Set("metadata", details["metadata"])
}

func ValidatePhoneNumberSignUp(c *gin.Context) {
	var details map[string]interface{}

	if err := c.BindJSON(&details); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	phoneNumber, ok := details["phoneNumber"].(string)

	if ok && helpers.IsValidPhoneNumber(phoneNumber) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
	}

	password, ok := details["password"].(string)

	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
	}

	err := helpers.IsValidPassword(password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
	}

	c.Set("email", phoneNumber)
	c.Set("password", password)
	c.Set("metadata", details["metadata"])
}
