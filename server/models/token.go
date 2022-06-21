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
	UserId   *string     `json:"id"`
	Metadata interface{} `json:"metadata"`
	jwt.StandardClaims
}

// todo => update to public - private key set up (RSA256)

func NewToken(tokenType string, user *User) (*string, error) {
	secret := getUserClaimSecret(tokenType)
	duration := getUserClaimDuration(tokenType)

	claim := &Token{
		UserId:   user.ID,
		Metadata: user.Metadata,
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

func DecodeToken(tokenType string, token *string) (*Token, error) {
	parsedToken, err := jwt.ParseWithClaims(
		*token,
		&Token{},
		func(token *jwt.Token) (interface{}, error) {
			return getUserClaimSecret(tokenType), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*Token)
	if !ok {
		return nil, errors.New("unable retrieve token")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("token is expired")
	}

	return claims, nil
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
