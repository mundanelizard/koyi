package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/handlers/middlewares"
	"github.com/mundanelizard/koyi/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

func extractIntentDetails(c *gin.Context) (string, string, error) {
	var intentId string
	var code string

	var err error

	if c.Request.Method == http.MethodPost {
		var body map[string]string
		err = c.BindJSON(&body)
		intentId = body["intentId"]
		code = body["code"]
	} else {
		intentId = c.Param("intentId")
		code = c.Param("code")
	}

	if len(intentId) == 0 {
		return "", "", errors.New("bad request: missing 'intentId'")
	} else if len(code) == 0 {
		return "", "", errors.New("bad request: missing 'code'")
	}

	return intentId, code, err
}

func verifyHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), config.AverageServerTimeout)
	defer cancel()

	intentId, intentCode, err := extractIntentDetails(c)

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	intent, err := models.FindIntent(ctx, intentId, intentCode)

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	// check if intent has expired
	if intent.IsExpired() {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	user, err := models.FindUserById(ctx, intent.UserId)

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{})
		return
	}

	if intent.Action == models.AccountVerificationIntent {
		err = user.Update(ctx, bson.M{"isEmailVerified": true})
	} else if intent.Action == models.AccountVerificationIntent {
		err = user.Update(ctx, bson.M{"isPhoneNumberVerified": true})
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}

	err = intent.Update(ctx, bson.M{"fulfilled": true})

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{})
}

func requestPhoneNumberVerificationHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), config.AverageServerTimeout)
	defer cancel()

	userId := c.GetString("userId")

	user, err := models.FindUserById(ctx, userId)

	if user.IsPhoneNumberVerified {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{})
		return
	}

	err = user.SendVerificationSms(ctx)

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{})
}

func requestEmailVerificationHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), config.AverageServerTimeout)
	defer cancel()

	userId := c.GetString("userId")

	user, err := models.FindUserById(ctx, userId)

	if user.IsEmailVerified {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{})
		return
	}

	err = user.SendVerificationMail(ctx)

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{})
}

func CreateVerificationRoutes(router *gin.RouterGroup) {
	// This route should be able to auto generate otp for admin application only, based on API Key
	group := router.Group("/auth/verify")

	group.POST("/phone-number", middlewares.AuthoriseUser, requestPhoneNumberVerificationHandler)
	group.POST("/email", middlewares.AuthoriseUser, requestEmailVerificationHandler)

	group.POST("/", verifyHandler)
	group.GET("/:intentId/:code", verifyHandler)
}
