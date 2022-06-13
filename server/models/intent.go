package models

import (
	"context"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/helpers"
)

const (
	intentsCollectionName = "intents"

	// Intents
	accountVerificationIntent = "verification"
)

type Intent struct {
	ID     string `json:"id"`
	UserId string `json:"userId"`

	Action string `json:"action"` // reset-password

	ActionCode string `json:"-"`
	ActionUrl  string `json:"-"`

	Fulfilled bool `json:"fulfilled"`
}

func (i *Intent) Create(ctx context.Context) error {
	i.ID = uuid.New().String()
	i.Fulfilled = false

	collection := helpers.GetCollection(config.UserDatabaseName, intentsCollectionName)
	_, err := collection.InsertOne(ctx, i)

	return err
}

func NewIntent(userId, action string, generateActionUrl func(intentId, actionCode string) string) *Intent {
	actionCode := helpers.RandomIntegers(6)
	intentId := uuid.New().String()

	intent := &Intent{
		ID:     intentId,
		UserId: userId,
		Action: action,
		// todo => change to template.
		ActionUrl:  generateActionUrl(intentId, actionCode),
		Fulfilled:  false,
		ActionCode: actionCode,
	}

	return intent
}

func getEmail(action string) (*string, *string) {
	text := ""
	html := ""
	return &text, &html
}

func getSms(action string) *string {
	text := ""
	return &text
}
