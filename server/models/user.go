package models

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"log"
	"time"
)

type User struct {
	Email       *string      `json:"email"`
	PhoneNumber *string      `json:"phoneNumber"`
	Password    *string      `json:"-" bson:"password"`
	Metadata    *interface{} `json:"metadata"`
	ID          *string      `json:"id"`
	IsDeleted   bool         `json:"deleted"`
	IsVerified  bool         `json:"isVerified"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

type UserClaim struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
	ID          *string `json:"id"`
	jwt.StandardClaims
}

func Count(interface{}) (int8, error) {
	return 0, nil
}

func (user *User) Exists() (bool, error) {
	var count int8
	var err error

	if user.Email != nil {
		count, err = Count(map[string]string{"email": *user.Email})
	} else if user.PhoneNumber != nil {
		count, err = Count(map[string]string{"email": *user.Email})
	} else {
		return false, errors.New("empty user object")
	}

	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (user *User) FillDefaults() {
	id := uuid.New().String()
	user.ID = &id
	user.IsDeleted = false
	user.IsVerified = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}

func (user *User) Create(ctx *context.Context) error {
	exists, err := user.Exists()

	if err != nil {
		return err
	}

	if exists {
		return errors.New("user already exits")
	}

	user.FillDefaults()

	return errors.New("test Error")
}

func (user *User) GenerateJWTs() (*string, *string, error) {
	accessToken, err := user.CreateClaim(config.AccessTokenDuration, config.AccessTokenSecretKey)
	refreshToken, err := user.CreateClaim(config.RefreshTokenDuration, config.RefreshTokenSecretKey)

	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (user *User) CreateClaim(duration time.Duration, secret string) (*string, error) {
	accessClaims := &UserClaim{
		ID:          user.ID,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		StandardClaims: jwt.StandardClaims{
			// Stays valid for 10 hours
			ExpiresAt: time.Now().Local().
				Add(duration).Unix(),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).
		SignedString([]byte(secret))

	return &accessToken, err
}
