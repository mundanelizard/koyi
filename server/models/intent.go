package models

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/helpers"
)

const (
	intentsCollectionName   = "intents"
	verifyPhoneNumberIntent = "verify:phone-number"
	verifyEmailIntent       = "verify:email"
)

var actionUrls = map[string]string{
	verifyEmailIntent:       "https://gmail.com/google.com",
	verifyPhoneNumberIntent: "https://email/intent/intent",
}

var actionSubject = map[string]string{}
var actionHTML = map[string]string{}
var actionText = map[string]string{}

type Intent struct {
	ID     string `json:"id"`
	UserId string `json:"userId"`

	Action string `json:"action"` // reset-password

	ActionCode string `json:"-"`
	ActionUrl  string `json:"-"`

	Fulfilled bool `json:"fulfilled"`
}

func CreateVerificationIntent(ctx *context.Context, user *User) (error) {
	actionCode := "123456"
	action, err := getVerificationIntentType(user)

	if err != nil {
		return err
	}

	intent := &Intent{
		UserId:     *user.ID,
		Action:     action,
		ActionCode: actionCode,
		ActionUrl:  actionUrls[action],
	}

	err = intent.Create(ctx)

	if err != nil {
		return err
	}

	return nil
}

func getVerificationIntentType(user *User) (string, error) {
	switch {
	case user.Email != nil:
		return verifyEmailIntent, nil
	case user.PhoneNumber != nil:
		return verifyPhoneNumberIntent, nil
	default:
		return "", errors.New("invalid verification intent type")
	}
}

func (i *Intent) Create(ctx *context.Context) error {
	i.ID = uuid.New().String()
	i.Fulfilled = false

	collection := helpers.GetCollection(config.UserDatabaseName, intentsCollectionName)
	_, err := collection.InsertOne(*ctx, i)

	return err
}