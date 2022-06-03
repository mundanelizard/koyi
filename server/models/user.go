package models

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/helpers"
	"log"
	"time"
)

const (
	userCollectionName = "users"
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

func Count(ctx *context.Context, filter interface{}) (int64, error) {
	collection := helpers.GetCollection(config.UserDatabaseName, userCollectionName)
	count, err := collection.CountDocuments(*ctx, filter)
	return count, err
}

func (user *User) Exists(ctx *context.Context) (bool, error) {
	var count int64
	var err error

	if user.Email != nil {
		count, err = Count(ctx, map[string]string{"email": *user.Email})
	} else if user.PhoneNumber != nil {
		count, err = Count(ctx, map[string]string{"email": *user.Email})
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
	exists, err := user.Exists(ctx)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("user already exits")
	}

	user.FillDefaults()
	collection := helpers.GetCollection(config.UserDatabaseName, userCollectionName)

	_, err = collection.InsertOne(*ctx, user)

	if err != nil {
		return err
	}

	return nil
}

func (user *User) PersistJWTs(accessToken, refreshToken *string) error {
	// store the JWT in the
	return nil
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
	claims := &UserClaim{
		ID:          user.ID,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		StandardClaims: jwt.StandardClaims{
			// Stays valid for 10 hours
			ExpiresAt: time.Now().Local().
				Add(duration).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(secret))

	return &token, err
}
