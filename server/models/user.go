package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const (
	userCollectionName        = "users"
	tokenClaimsCollectionName = "user-tokens"
)

type User struct {
	Email                 *string      `json:"email"`
	IsEmailVerified       bool         `json:"isEmailVerified"`
	PhoneNumber           *PhoneNumber `json:"phoneNumber"`
	IsPhoneNumberVerified bool         `json:"isPhoneNumberVerified"`
	Password              *string      `json:"-" bson:"password"`
	Metadata              *interface{} `json:"metadata"`
	ID                    *string      `json:"id"`
	IsDeleted             bool         `json:"deleted"`
	CreatedAt             time.Time    `json:"createdAt"`
	UpdatedAt             time.Time    `json:"updatedAt"`
}

func CountUser(ctx context.Context, filter interface{}) (int64, error) {
	collection := helpers.GetCollection(config.UserDatabaseName, userCollectionName)
	count, err := collection.CountDocuments(ctx, filter)
	return count, err
}

func (user *User) Create(ctx context.Context) error {
	exists, err := user.exists(ctx)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("user already exits")
	}

	user.fillDefaults()
	collection := helpers.GetCollection(config.UserDatabaseName, userCollectionName)

	_, err = collection.InsertOne(ctx, user)

	if err != nil {
		return err
	}

	go user.createHistory(ctx)

	return nil
}

func FindUser(ctx context.Context, filter bson.M) (*User, error) {

	return nil, nil
}

func (user *User) CreateClaims(ctx context.Context, deviceId string) (*TokenClaims, error) {
	accessClaim, accessToken, err := NewUserClaim("access", user)
	if err != nil {
		return nil, err
	}

	refreshClaim, refreshToken, err := NewUserClaim("refresh", user)
	if err != nil {
		return nil, err
	}

	claims := NewTokenClaim(accessToken, refreshToken, refreshClaim, accessClaim, &deviceId)

	err = claims.Create(ctx)

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (user *User) SendVerificationMessage(ctx context.Context) error {
	intent := NewIntent(
		*user.ID,
		accountVerificationIntent,
		func(intentId, actionId string) string {
			return fmt.Sprintf(
				config.ServerDomain+"/v1/auth/signup/verify/%s/%s",
				intentId, actionId)
		},
	)

	err := intent.Create(ctx)

	if err != nil {
		return err
	}

	data := &templateData{Intent: intent, User: user}
	var m helpers.Sendable

	if user.Email != nil {
		m, err = getEmail(data)
	} else if user.PhoneNumber != nil {
		m, err = getSms(data)
	}

	if err != nil {
		return err
	}

	return m.Send()
}

// INTERNALS

// fillDefaults sets the User.IsDeleted, User.IsVerified, User.CreatedAt, User.UpdatedAt
// User.ID fields on the struct. If the User.ID field is already present, don't provide any default
// for the User struct fields.
func (user *User) fillDefaults() {
	if user.ID != nil {
		return
	}

	id := uuid.New().String()
	user.ID = &id
	user.IsDeleted = false

	user.IsEmailVerified = false
	user.IsPhoneNumberVerified = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}

// createHistory creates an EmailHistory, PasswordHistory and PhoneNumberHistory for a User.
func (user *User) createHistory(ctx context.Context) {
	eh := NewEmailHistory(user.ID, user.Email)
	go eh.Create(ctx)

	ph := NewPasswordHistory(user.ID, user.Password)
	go ph.Create(ctx)

	pnh := NewPhoneNumberHistory(user.ID, user.PhoneNumber)
	go pnh.Create(ctx)
}

// exists checks if a user exists in the database.
func (user *User) exists(ctx context.Context) (bool, error) {
	var count int64
	var err error

	if user.Email != nil {
		count, err = CountUser(ctx, map[string]string{"email": *user.Email})
	} else if user.PhoneNumber != nil {
		count, err = CountUser(ctx,
			bson.M{
				"phoneNumber.countryCode":      user.PhoneNumber.CountryCode,
				"phoneNumber.subscriberNumber": user.PhoneNumber.SubscriberNumber,
			})
	} else {
		return false, errors.New("something terribly wrong happened: the user doesn't have an email or phone number")
	}

	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}
