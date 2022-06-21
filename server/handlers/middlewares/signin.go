package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/mundanelizard/koyi/server/models"
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

	if !ok || ValidateEmail(email) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	password, ok := details["password"].(string)

	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	err := ValidatePassword(password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
	}

	c.Set("email", email)
	c.Set("password", password)
	c.Next()
}

func ValidatePhoneNumberSignIn(c *gin.Context) {
	var details map[string]interface{}

	if err := c.BindJSON(&details); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	sn, _ := details["subscriberNumber"].(string)
	cc, _ := details["countryCode"].(string)

	phoneNumber := &models.PhoneNumber{
		SubscriberNumber: sn,
		CountryCode:      cc,
	}

	if !phoneNumber.IsValid() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	password, ok := details["password"].(string)

	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	err := ValidatePassword(password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.Set("phoneNumber", phoneNumber)
	c.Set("password", password)
	c.Next()
}
