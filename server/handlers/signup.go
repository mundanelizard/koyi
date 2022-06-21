package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/handlers/middlewares"
	"github.com/mundanelizard/koyi/server/models"
	"log"
	"net/http"
)

func emailSignUpHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), config.AverageServerTimeout)
	defer cancel()

	email := c.GetString("email")
	password := c.GetString("password")
	metadata, _ := c.Get("metadata")

	user := &models.User{
		Email:    &email,
		Password: &password,
		Metadata: &metadata,
	}

	err := user.Create(ctx)

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	err = user.SendVerificationMail(ctx)

	if err != nil {
		log.Println("SEND-VERIFICATION-EMAIL-ERROR: ", err)
	}

	// todo => remove the duplications
	if config.CreateTokenOnSignUp {
		device := models.ExtractAndCreateDevice(ctx, c.Request, *user.ID)
		AbortGinWithAuth(c, ctx, user, device.ID)
	} else {
		go models.ExtractAndCreateDevice(ctx, c.Request, *user.ID)
		c.AbortWithStatusJSON(http.StatusCreated, user)
	}
}

func phoneNumberSignUpHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), config.AverageServerTimeout)
	defer cancel()

	phoneNumber, _ := c.Get("phoneNumber")
	password := c.GetString("password")
	metadata, _ := c.Get("metadata")

	user := &models.User{
		PhoneNumber: phoneNumber.(*models.PhoneNumber),
		Password:    &password,
		Metadata:    &metadata,
	}

	err := user.Create(ctx)

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	err = user.SendVerificationSms(ctx)

	if err != nil {
		log.Println("SEND-VERIFICATION-SMS-ERROR: ", err)
	}

	// todo => remove the duplications
	if config.CreateTokenOnSignUp {
		device := models.ExtractAndCreateDevice(ctx, c.Request, *user.ID)
		AbortGinWithAuth(c, ctx, user, device.ID)
	} else {
		go models.ExtractAndCreateDevice(ctx, c.Request, *user.ID)
		c.AbortWithStatusJSON(http.StatusCreated, user)
	}
}

func CreateSignUpRoutes(router *gin.RouterGroup) {
	// Add the ability to assign roles to users.
	// If the user is coming from 'web' you can restrict his roles.
	group := router.Group("/auth/signup")

	group.POST("/email", middlewares.EmailSignUpValidator, emailSignUpHandler)
	group.POST("/phone", middlewares.PhoneNumberSignUpValidator, phoneNumberSignUpHandler)
}
