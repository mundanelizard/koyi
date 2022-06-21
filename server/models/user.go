package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/services"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

const (
	userCollectionName        = "users"
	tokenClaimsCollectionName = "user-tokens"
)

type User struct {
	Email                 *string      `json:"email" bson:"email"`
	IsEmailVerified       bool         `json:"isEmailVerified" bson:"isEmailVerified"`
	PhoneNumber           *PhoneNumber `json:"phoneNumber" bson:"phoneNumber"`
	IsPhoneNumberVerified bool         `json:"isPhoneNumberVerified" bson:"isPhoneNumberVerified"`
	Password              *string      `json:"-" bson:"password"`
	Metadata              *interface{} `json:"metadata" bson:"metadata"`
	ID                    *string      `json:"id" bson:"id"`
	IsDeleted             bool         `json:"deleted" bson:"isDeleted"`
	CreatedAt             time.Time    `json:"createdAt" bson:"createdAt"`
	UpdatedAt             time.Time    `json:"updatedAt" bson:"updatedAt"`
}

func CountUser(ctx context.Context, filter interface{}) (int64, error) {
	collection := services.GetCollection(config.UserDatabaseName, userCollectionName)
	count, err := collection.CountDocuments(ctx, filter)
	return count, err
}

func (user *User) Create(ctx context.Context) error {
	exists, err := user.exists(ctx)

	if exists || err != nil {
		return errors.New("user already exits")
	}

	if user.Password == nil {
		return errors.New("user password error")
	}

	user.fillDefaults()
	user.Password = hashPassword(*user.Password)
	collection := services.GetCollection(config.UserDatabaseName, userCollectionName)

	_, err = collection.InsertOne(ctx, user)

	if err != nil {
		return err
	}

	go user.createHistory(ctx)

	return nil
}

func FindUser(ctx context.Context, filter bson.M) (*User, error) {
	var user User

	collection := services.GetCollection(config.UserDatabaseName, intentsCollectionName)
	err := collection.FindOne(ctx, filter).Decode(&user)

	return &user, err
}

func FindUserById(ctx context.Context, userId string) (*User, error) {
	var user User

	collection := services.GetCollection(config.UserDatabaseName, intentsCollectionName)
	err := collection.FindOne(ctx, bson.M{"id": userId}).Decode(&user)

	return &user, err
}

func (user *User) Update(ctx context.Context, m bson.M) error {
	collection := services.GetCollection(config.UserDatabaseName, intentsCollectionName)
	return collection.FindOneAndUpdate(ctx, bson.M{"id": user.ID}, m).Decode(user)
}

func (user *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password))

	if err != nil {
		return false
	}

	return true
}

func (user *User) CreateClaims(ctx context.Context, deviceId string) (*TokenClaims, error) {
	at, err := NewToken(accessToken, user)
	if err != nil {
		return nil, err
	}

	rt, err := NewToken(refreshToken, user)
	if err != nil {
		return nil, err
	}

	claims := NewClaim(at, rt, user.ID, &deviceId)

	err = claims.Create(ctx)

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (user *User) SendVerificationMail(ctx context.Context) error {
	intent := NewIntent(
		*user.ID,
		EmailVerificationIntent,
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
	m, err := getEmail(data)

	if err != nil {
		return err
	}

	return m.Send()
}

func (user *User) SendVerificationSms(ctx context.Context) error {
	intent := NewIntent(
		*user.ID,
		PhoneNumberVerificationIntent,
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
	m, err := getSms(data)

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
	ph := NewHistory(*user.ID, passwordFieldName, *user.Password)
	ph.Create(ctx)

	if user.Email != nil {
		eh := NewHistory(*user.ID, emailFieldName, *user.Email)
		eh.Create(ctx)
	}

	if user.PhoneNumber != nil {
		pnh := NewHistory(*user.ID, phoneNumberFieldName, user.PhoneNumber)
		pnh.Create(ctx)
	}
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

func hashPassword(password string) *string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Fatal(err)
	}

	result := string(bytes)

	return &result
}
