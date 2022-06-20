package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/helpers"
	"github.com/mundanelizard/koyi/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func EmailSignInHandler(c *gin.Context) {
	email := c.GetString("email")
	password := c.GetString("password")

	user, err := models.FindUser(c, bson.M{
		"email": email,
	})

	if err != nil {
		log.Println(err)
		// todo => return a 404 error.
	}

	SignInHandler(c, user, password)
}

func PhoneNumberSignInHandler(c *gin.Context) {
	countryCode := c.GetString("countryCode")
	subscriberNumber := c.GetString("subscriberNumber")
	password := c.GetString("password")

	user, err := models.FindUser(c, bson.M{
		"phoneNumber.subscriberNumber": subscriberNumber,
		"phoneNumber.countryCode":      countryCode,
	})

	if err != nil {
		log.Println(err)
		// todo => return a 404 error.
	}

	SignInHandler(c, user, password)
}

func SignInHandler(c *gin.Context, user *models.User, password string) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 100)
	defer cancel()

	ok := helpers.VerifyHash(&password, user.Password)

	if !ok {
		log.Println("wrong password")
		// todo => return a 400 error
	}

	device := models.ExtractAndCreateDevice(c.Request, user.ID)
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

	AbortGinWithAuth(c, ctx, user, *device.ID)
}
