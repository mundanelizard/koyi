package models

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"time"
)

const (
	accessToken  = "access-token"
	refreshToken = "refresh-token"
)

type Token struct {
	Email       *string      `json:"email"`
	PhoneNumber *PhoneNumber `json:"phoneNumber"`
	ID          *string      `json:"id"`
	jwt.StandardClaims
}

func NewToken(tokenType string, user *User) (*string, error) {
	secret := getUserClaimSecret(tokenType)
	duration := getUserClaimDuration(tokenType)

	claim := &Token{
		ID:          user.ID,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		StandardClaims: jwt.StandardClaims{
			Issuer:   config.JWTIssuerName,
			IssuedAt: time.Now().Unix(),
			Audience: config.JWTAudience,
			Id:       uuid.New().String(),
			// todo => NotBefore "nbf"
			// todo => Subject   "sub"
			// Stays valid for 10 hours
			ExpiresAt: time.Now().Local().
				Add(duration).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).
		SignedString([]byte(secret))

	if err != nil {
		return &token, errors.New(fmt.Sprintf("CREATE-JWT-ERROR: %s", err))
	}

	return &token, err
}

func getUserClaimSecret(claimType string) string {
	if claimType == refreshToken {
		return config.RefreshTokenSecretKey
	}

	return config.AccessTokenSecretKey
}

func getUserClaimDuration(claimType string) time.Duration {
	if claimType == refreshToken {
		return config.RefreshTokenDuration
	}

	return config.AccessTokenDuration
}
