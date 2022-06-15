package models

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/helpers"
	"html/template"
)

const (
	intentsCollectionName = "intents"

	// Intents
	accountVerificationIntent = "verification"
)

var (
	htmlEmailVerificationTemplate = template.Must(template.ParseFiles("./templates/verification.html"))
	textEmailVerificationTemplate = template.Must(template.ParseFiles("./templates/verification.html"))
	smsVerificationTemplate       = template.Must(template.ParseFiles("./templates/verification.html"))
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

type templateData struct {
	Intent *Intent
	User   *User
}

func getEmail(data *templateData) (helpers.Sendable, error) {
	var textBuffer *bytes.Buffer
	var htmlBuffer *bytes.Buffer
	var subject string
	var err error

	if data.Intent.Action == accountVerificationIntent {
		subject = "Verification Email"
		err = htmlEmailVerificationTemplate.Execute(htmlBuffer, data)
		err = textEmailVerificationTemplate.Execute(textBuffer, data)
	}

	if err != nil {
		return nil, err
	}

	text := textBuffer.String()
	html := htmlBuffer.String()

	return helpers.NewMail(*data.User.Email, subject, &text, &html), err
}

func getSms(data *templateData) (helpers.Sendable, error) {
	var buffer *bytes.Buffer
	var err error

	if data.Intent.Action == accountVerificationIntent {
		err = smsVerificationTemplate.Execute(buffer, data)
	}

	if err != nil {
		return nil, err
	}

	text := buffer.String()

	return helpers.NewSms(&text), err
}
