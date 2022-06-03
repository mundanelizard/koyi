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

	metadata, ok := c.Get("metadata")

	if ok {
		user.Metadata = &metadata
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

	metadata, ok := c.Get("metadata")

	if ok {
		user.Metadata = &metadata
	}

	CreateUserHandler(c, user)
}

func CreateUserHandler(c *gin.Context, user *models.User) {
	device := helpers.ExtractDeviceDetailsFromContext(c)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 100)
	defer cancel()

	err := user.Create(&ctx)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	device.UserId = user.ID
	err = device.Create(&ctx)

	if err != nil {
		log.Println("CREATE-DEVICE-ERROR: ", err)
	}

	accessToken, refreshToken, err := user.GenerateTokensAndPersistClaims(&ctx, device)

	if err != nil {
		log.Println("GENERATE-JWT-ERROR: ", err)
	}

	response := map[string]interface{}{
		"user":         user,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	c.SetCookie(
		"authentication",
		*refreshToken,
		config.RefreshTokenCookieMaxAge,
		"/v1/who/refresh",
		config.ServerDomain,
		config.IsProduction,
		true)

	c.AbortWithStatusJSON(http.StatusCreated, response)
}
