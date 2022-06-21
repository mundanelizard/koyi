package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/models"
	"log"
	"net/http"
)

func AbortGinWithAuth(c *gin.Context, ctx context.Context, user *models.User, deviceId string) {
	claims, err := user.CreateTokens(ctx, deviceId)

	if err != nil {
		log.Println(err)
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
