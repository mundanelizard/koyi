package middlewares

import (
	"github.com/gin-gonic/gin"
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

	token := strings.Split(header, "JWT ") // Bearer || JWT || ACCESS-KEY

	if len(token) < 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}

	t, err := models.DecodeToken("access-token", &token[0])

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}

	claim, err := models.FindAccessClaim(ctx, t.UserId, token)

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

	c.Set("token", t)
}
