package models

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"time"
)

type UserClaim struct {
	Email       *string      `json:"email"`
	PhoneNumber *PhoneNumber `json:"phoneNumber"`
	ID          *string      `json:"id"`
	jwt.StandardClaims
}

func NewUserClaim(tokenType string, user *User) (*UserClaim, *string, error) {
	secret := getUserClaimSecret(tokenType)
	duration := getUserClaimDuration(tokenType)

	claims := &UserClaim{
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

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(secret))

	return claims, &token, err
}

func getUserClaimSecret(claimType string) string {
	if claimType == "refresh" {
		return config.RefreshTokenSecretKey
	}

	return config.AccessTokenSecretKey
}

func getUserClaimDuration(claimType string) time.Duration {
	if claimType == "refresh" {
		return config.RefreshTokenDuration
	}

	return config.AccessTokenDuration
}