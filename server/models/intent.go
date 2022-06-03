package models

import (
	"context"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/helpers"
)

const (
	intentsCollectionName   = "intents"
	VerifyPhoneNumberIntent = "verify:phone-number"
	VerifyEmailIntent       = "verify:email"
)

var actionUrls = map[string]string{
	VerifyEmailIntent:       "https://gmail.com/google.com",
	VerifyPhoneNumberIntent: "https://email/intent/intent",
}

type Intent struct {
	ID     string `json:"id"`
	UserId string `json:"userId"`

	Action string `json:"action"` // reset-password

	ActionCode string `json:"-"`
	ActionUrl  string `json:"-"`

	Fulfilled bool `json:"fulfilled"`
}

func CreateIntent(ctx *context.Context, user *User, action string) (*Intent, error) {
	actionCode := "123456"
	var err error

	if err != nil {
		return nil, err
	}

	intent := &Intent{
		UserId:     *user.ID,
		Action:     action,
		ActionCode: actionCode,
		ActionUrl:  actionUrls[action],
	}

	return intent, intent.Create(ctx)
}

func (i *Intent) Create(ctx *context.Context) error {
	i.ID = uuid.New().String()
	i.Fulfilled = false

	collection := helpers.GetCollection(config.UserDatabaseName, intentsCollectionName)
	_, err := collection.InsertOne(*ctx, i)

	return err
}
