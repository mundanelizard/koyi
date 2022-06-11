package helpers

import (
	"errors"
	"net/mail"
	"unicode"
)

// ValidateEmail validates email using a wrapper around mail.ParseAddress and only returns the error.
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
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

// IsValidSubscriberNumber validates subscriber ar https://en.wikipedia.org/wiki/E.164
func IsValidSubscriberNumber(sn string) bool {
	length := len(sn)
	return length > 3 && length <= 13
}

func IsValidCountryCode(cc string) bool {
	length := len(cc)
	return length > 1 && length <= 3
}
