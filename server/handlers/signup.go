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

	go user.SendVerificationMessage(ctx)
	go device.Create(ctx)

	AbortGinWithAuth(c, ctx, user, *device.ID)
}

// CreateAuthenticationRoutes handles all the authentication request made in the application.
func CreateAuthenticationRoutes(router *gin.RouterGroup) {
	group := router.Group("/auth/signup")

	group.POST("/email", middlewares.ValidateEmailSignUp, emailSignUpHandler)
	group.POST("/phone", middlewares.ValidatePhoneNumberSignUp, phoneNumberSignUpHandler)
}
