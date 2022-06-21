package models

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/services"
	"go.mongodb.org/mongo-driver/bson"
	"html/template"
	"time"
)

const (
	intentsCollectionName = "intents"

	// -- intent types
	AccountVerificationIntent     = "verification"
	EmailVerificationIntent       = "verification:email"
	PhoneNumberVerificationIntent = "verification:phone-number"
)

var (
	htmlEmailVerificationTemplate = template.Must(template.ParseFiles(config.HTMLEmailVerificationTemplatePath))
	textEmailVerificationTemplate = template.Must(template.ParseFiles(config.TextEmailVerificationTemplatePath))
	smsVerificationTemplate       = template.Must(template.ParseFiles(config.PhoneNumberVerificationTemplatePath))
)

type Intent struct {
	ID     string `json:"id" bson:"id"`
	UserId string `json:"userId" bson:"userId"`

	Action string `json:"action" bson:"action"` // reset-password

	ActionCode string `json:"-" bson:"actionCode"`
	ActionUrl  string `json:"-" bson:"actionUrl"`

	Fulfilled bool      `json:"fulfilled" bson:"fulfilled"`
	ExpireAt  time.Time `json:"expireAt" bson:"expireAt"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func (i *Intent) Create(ctx context.Context) error {
	i.ID = uuid.New().String()
	i.Fulfilled = false

	collection := services.GetCollection(config.UserDatabaseName, intentsCollectionName)
	_, err := collection.InsertOne(ctx, i)

	return err
}

func NewIntent(userId, action string, generateActionUrl func(intentId, actionCode string) string) *Intent {
	actionCode := services.RandomIntegers(6)
	intentId := uuid.New().String()

	intent := &Intent{
		ID:     intentId,
		UserId: userId,
		Action: action,
		// todo => change to template.
		ActionUrl:  generateActionUrl(intentId, actionCode),
		Fulfilled:  false,
		ActionCode: actionCode,
		ExpireAt:   time.Now().Local().Add(config.IntentDuration),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return intent
}

func FindIntent(ctx context.Context, intentId, intentCode string) (*Intent, error) {
	var intent Intent

	collection := services.GetCollection(config.UserDatabaseName, intentsCollectionName)
	err := collection.FindOne(ctx, bson.M{"id": intentId, "actionCode": intentCode}).Decode(&intent)

	return &intent, err
}

func (i *Intent) IsExpired() bool {
	return i.ExpireAt.Before(time.Now())
}

func (i *Intent) Update(ctx context.Context, m bson.M) error {
	collection := services.GetCollection(config.UserDatabaseName, intentsCollectionName)
	return collection.FindOneAndUpdate(ctx, bson.M{"id": i.ID}, m).Decode(i)
}

type templateData struct {
	Intent *Intent
	User   *User
}

func getEmail(data *templateData) (services.Sendable, error) {
	var textBuffer *bytes.Buffer
	var htmlBuffer *bytes.Buffer
	var subject string
	var err error

	if data.Intent.Action == AccountVerificationIntent {
		subject = "Verification Email"
		err = htmlEmailVerificationTemplate.Execute(htmlBuffer, data)
		err = textEmailVerificationTemplate.Execute(textBuffer, data)
	}

	if err != nil {
		return nil, err
	}

	text := textBuffer.String()
	html := htmlBuffer.String()

	return services.NewMail(*data.User.Email, subject, &text, &html), err
}

func getSms(data *templateData) (services.Sendable, error) {
	var buffer *bytes.Buffer
	var err error

	if data.Intent.Action == AccountVerificationIntent {
		err = smsVerificationTemplate.Execute(buffer, data)
	}

	if err != nil {
		return nil, err
	}

	text := buffer.String()

	return services.NewSms(&text), err
}
