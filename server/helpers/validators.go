package helpers

import (
	"errors"
	"net/mail"
	"unicode"
)

// IsValidEmail validates email using a wrapper around mail.ParseAddress and only returns the error.
func IsValidEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

// IsValidPassword validates a password based on the length, special characters and casing.
func IsValidPassword(password string) error {
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

// IsValidPhoneNumber validates phone number to match a simplified version of Nigeria phone number spec.
func IsValidPhoneNumber(phoneNumber string) error {
	firstNumber := phoneNumber[0]

	// this only works for west african phone numbers.
	if firstNumber != '0' {
		return errors.New("invalid phone number")
	}

	if len(phoneNumber) != 11 {
		return errors.New("invalid phone number")
	}

	return nil
}
