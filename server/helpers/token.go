package helpers

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashString(password string) *string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Fatal(err)
	}

	result := string(bytes)

	return &result
}

func VerifyHash(password, hash *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*hash), []byte(*password))

	if err != nil {
		return false
	}

	return true
}
