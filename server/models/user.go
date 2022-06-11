package models

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/pkg/email"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

const (
	userCollectionName        = "users"
	tokenClaimsCollectionName = "user-tokens"
)

type User struct {
	Email       *string      `json:"email"`
	PhoneNumber *PhoneNumber `json:"phoneNumber"`
	Password    *string      `json:"-" bson:"password"`
	Metadata    *interface{} `json:"metadata"`
	ID          *string      `json:"id"`
	IsDeleted   bool         `json:"deleted"`
	IsVerified  bool         `json:"isVerified"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
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

func (user *User) SendVerificationMessage(ctx context.Context) {
	intent, err := CreateIntent(ctx, user, getVerificationIntentType(user))

	if err != nil {
		log.Println("CREATE-VERIFICATION-INTENT-ERROR: ", err)
	}

	err = user.sendMessage(ctx, intent)

	if err != nil {
		log.Println("SEND-VERIFICATION-EMAIL-ERROR: ", err)
	}
}

// INTERNALS

// sendMessage sends an email or sms to a user.
func (user *User) sendMessage(ctx context.Context, intent *Intent) error {
	defer ctx.Done()

	switch intent.Action {
	case VerifyEmailIntent:
		return user.sendEmail(intent)
	case VerifyPhoneNumberIntent:
		return user.sendEmail(intent)
	}

	return errors.New("unable to find email intent")
}

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
	user.IsVerified = false
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

// sendEmail sends an email to the user (right now using amazon ses)
func (user *User) sendEmail(intent *Intent) error {
	// todo => compile template file
	subject, text, html := GetEmailDetails(intent.Action)

	e := &email.Email{
		Subject:  subject,
		BodyText: text,
		BodyHTML: html,
	}

	return e.Send()
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
				"$and": bson.M{
					"phoneNumber.countryCode":      user.PhoneNumber.CountryCode,
					"phoneNumber.subscriberNumber": user.PhoneNumber.SubscriberNumber,
				},
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
