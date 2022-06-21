package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/handlers/middlewares"
	"github.com/mundanelizard/koyi/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

func emailSignInHandler(c *gin.Context) {
	email := c.GetString("email")
	password := c.GetString("password")

	user, err := models.FindUser(c, bson.M{
		"email": email,
	})

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{})
		return
	}

	signIn(c, user, password)
}

func phoneNumberSignInHandler(c *gin.Context) {
	phoneNumber, _ := c.Get("phoneNumber")
	password := c.GetString("password")

	user, err := models.FindUser(c, bson.M{"phoneNumber": phoneNumber.(*models.PhoneNumber)})

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{})
		return
	}

	signIn(c, user, password)
}

func signIn(c *gin.Context, user *models.User, password string) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 100)
	defer cancel()

	ok := user.VerifyPassword(password)

	if !ok {
		log.Println("wrong password")
		// todo => return a 400 error
	}

	device := models.ExtractDevice(c.Request)
	exists, err := device.Exists(ctx)

	if err != nil {
		log.Println("CHECK-DEVICE-EXISTS-ERROR:", err)
	}

	if !exists && config.ValidateNewDevice {
		// create an intent to validate new device.
		// send an email for the intent.
		// send user to the validation route.
		log.Println("error")
		// todo => send an email giving the user the ability to invalidate tokens.
	}

	AbortGinWithAuth(c, ctx, user, device.ID)
}

func CreateSignInRoutes(router *gin.RouterGroup) {
	group := router.Group("/auth/signin")

	group.POST("/phone-number", middlewares.ValidatePhoneNumberSignIn, phoneNumberSignInHandler)
	group.POST("/email", middlewares.ValidateEmailSignIn, emailSignInHandler)
}
