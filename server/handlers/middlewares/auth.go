package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/models"
	"log"
	"net/http"
	"strings"
)

func AuthoriseUser(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")

	if len(header) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}

	segments := strings.Split(header, "JWT ") // Bearer || JWT || ACCESS-KEY

	if len(segments) < 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}

	t, err := models.DecodeClaim("access-token", &segments[1])

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.AverageServerTimeout)
	defer cancel()

	claim, err := models.FindClaim(ctx, *t.UserId, "access-token", segments[1])

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}

	if claim == nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}

	if *claim.DeviceId != c.Request.Header.Get("device-id") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}

	c.Set("claim", claim)
	c.Next()
}
