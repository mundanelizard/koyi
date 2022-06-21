package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mundanelizard/koyi/server/models"
	"log"
	"net/http"
	"net/mail"
	"unicode"
)

func EmailSignUpValidator(c *gin.Context) {
	var details map[string]interface{}

	if err := c.BindJSON(&details); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	email, ok := details["email"].(string)

	if !ok || ValidateEmail(email) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	password, ok := details["password"].(string)

	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	err := ValidatePassword(password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
	}

	c.Set("email", email)
	c.Set("password", password)
	c.Set("metadata", details["metadata"])
}

func PhoneNumberSignUpValidator(c *gin.Context) {
	var details map[string]interface{}

	if err := c.BindJSON(&details); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	phoneNumber := &models.PhoneNumber{
		SubscriberNumber: details["countryCode"].(string),
		CountryCode:      details["subscriberNumber"].(string),
	}

	if !phoneNumber.IsValid() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	password, ok := details["password"].(string)

	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	err := ValidatePassword(password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.Set("phoneNumber", &phoneNumber)
	c.Set("password", password)
	c.Set("metadata", details["metadata"])
}

// ValidatePassword validates a password based on the length, special characters and casing.
func ValidatePassword(password string) error {
	var count int
	var containsNumber, containsUppercase, containsSpecialCharacter bool

	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			containsNumber = true
		case unicode.IsUpper(c):
			containsUppercase = true
			count++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			containsSpecialCharacter = true
		case unicode.IsLetter(c) || c == ' ':
			count++
		default:
			return errors.New("invalid password")
		}
	}

	if !containsNumber {
		return errors.New("password doesn't contain a number")
	}

	if !containsUppercase {
		return errors.New("password doesn't contain an uppercase letter")
	}

	if !containsSpecialCharacter {
		return errors.New("password doesn't contain special character")
	}

	if count < 8 {
		return errors.New("password is less than 8 characters")
	}

	return nil
}

// ValidateEmail validates email using a wrapper around mail.ParseAddress and only returns the error.
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}
