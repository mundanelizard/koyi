package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/handlers/middlewares"
	"github.com/mundanelizard/koyi/server/helpers"
	"github.com/mundanelizard/koyi/server/models"
	"log"
	"net/http"
)

func emailSignUpHandler(c *gin.Context) {
	email := c.GetString("email")
	password := c.GetString("password")

	user := &models.User{
		Email:    &email,
		Password: helpers.HashString(password),
	}

	createUserHandler(c, user)
}

func phoneNumberSignUpHandler(c *gin.Context) {
	countryCode := c.GetString("countryCode")
	subscriberNumber := c.GetString("subscriberNumber")
	password := c.GetString("password")

	user := &models.User{
		PhoneNumber: &models.PhoneNumber{
			CountryCode:      countryCode,
			SubscriberNumber: subscriberNumber,
		},
		Password: helpers.HashString(password),
	}

	createUserHandler(c, user)
}

func createUserHandler(c *gin.Context, user *models.User) {
	ctx, cancel := context.WithTimeout(context.Background(), config.AverageServerTimeout)
	defer cancel()

	metadata, ok := c.Get("metadata")

	if ok {
		user.Metadata = &metadata
	}

	err := user.Create(ctx)

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	device := models.ExtractDevice(c.Request, user.ID)

	AbortGinWithAuth(c, ctx, user, *device.ID)

	err = user.SendVerificationMessage(ctx)

	if err != nil {
		log.Println("SEND-VERIFICATION-EMAIL-ERROR: ", err)
	}

	err = device.Create(ctx)

	if err != nil {
		log.Println(err)
	}
}

func CreateSignUpRoutes(router *gin.RouterGroup) {
	group := router.Group("/auth/signup")

	group.POST("/email", middlewares.ValidateEmailSignUp, emailSignUpHandler)
	group.POST("/phone", middlewares.ValidatePhoneNumberSignUp, phoneNumberSignUpHandler)
}
