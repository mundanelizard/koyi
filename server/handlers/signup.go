package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/helpers"
	"github.com/mundanelizard/koyi/server/models"
	"log"
	"net/http"
)

func EmailSignUpHandler(c *gin.Context) {
	email := c.GetString("email")
	password := c.GetString("password")

	user := &models.User{
		Email:    &email,
		Password: helpers.HashString(password),
	}

	CreateUserHandler(c, user)
}

func PhoneNumberSignUpHandler(c *gin.Context) {
	phoneNumber := c.GetString("phoneNumber")
	password := c.GetString("password")

	user := &models.User{
		PhoneNumber: &phoneNumber,
		Password:    helpers.HashString(password),
	}

	CreateUserHandler(c, user)
}

func CreateUserHandler(c *gin.Context, user *models.User) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 100)
	defer cancel()

	metadata, ok := c.Get("metadata")

	if ok {
		user.Metadata = &metadata
	}

	err := user.Create(&ctx)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	device := helpers.ExtractDeviceDetailsFromContext(c)
	device.UserId = user.ID
	err = device.Create(&ctx)

	if err != nil {
		log.Println("CREATE-DEVICE-ERROR: ", err)
	}

	intent, err := models.CreateIntent(&ctx, user, getVerificationIntentType(user))

	if err != nil {
		log.Println("CREATE-VERIFICATION-INTENT-ERROR: ", err)
	}

	go user.SendEmail(intent)

	claims, err := user.CreateClaims(&ctx, device)

	if err != nil {
		log.Println("GENERATE--JWT-ERROR: ", err)
	}

	response := map[string]interface{}{
		"user":    user,
		"token":   claims.AccessToken,
		"success": true,
	}

	c.SetCookie(
		"authentication",
		*claims.RefreshToken,
		config.RefreshTokenCookieMaxAge,
		"/v1/who/refresh",
		config.ServerDomain,
		config.IsProduction,
		true)

	c.AbortWithStatusJSON(http.StatusCreated, response)
}

func getVerificationIntentType(user *models.User) string {
	switch {
	case user.Email != nil:
		return models.VerifyEmailIntent
	case user.PhoneNumber != nil:
		return models.VerifyPhoneNumberIntent
	}

	return models.VerifyEmailIntent
}
