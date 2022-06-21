package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/mundanelizard/koyi/server/services"
	"log"
	"net/http"
)

func ValidateEmailSignIn(c *gin.Context) {
	var details map[string]interface{}

	if err := c.BindJSON(&details); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	email, ok := details["email"].(string)

	if ok && services.ValidateEmail(email) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
	}

	password, ok := details["password"].(string)

	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
	}

	err := services.ValidatePassword(password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
	}

	c.Set("email", email)
	c.Set("password", password)
}

func ValidatePhoneNumberSignIn(c *gin.Context) {
	var details map[string]interface{}

	if err := c.BindJSON(&details); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	countryCode, ok := details["countryCode"].(string)

	if ok && services.IsValidCountryCode(countryCode) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
	}

	subscriberNumber, ok := details["subscriberNumber"].(string)

	if ok && services.IsValidSubscriberNumber(subscriberNumber) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
	}

	password, ok := details["password"].(string)

	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
	}

	err := services.ValidatePassword(password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
	}

	c.Set("subscriberNumber", subscriberNumber)
	c.Set("countryCode", countryCode)
	c.Set("password", password)
}
