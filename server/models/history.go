package models

import (
	"context"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/helpers"
	"log"
	"time"
)

const (
	passwordHistoryCollectionName    = "password-history"
	phoneNumberHistoryCollectionName = "password-history"
	emailHistoryCollectionName       = "password-history"
)

type PasswordHistory struct {
	ID        *string   `json:"id"`
	UserId    *string   `json:"userId"`
	Password  *string   `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type EmailHistory struct {
	ID        *string   `json:"id"`
	UserId    *string   `json:"userId"`
	Email     *string   `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type PhoneNumberHistory struct {
	ID          *string      `json:"id"`
	UserId      *string      `json:"userId"`
	PhoneNumber *PhoneNumber `json:"phoneNumber"`
	CreatedAt   time.Time    `json:"createdAt"`
}

func (ph *PasswordHistory) FillDefaults() {
	id := uuid.New().String()
	ph.ID = &id
	ph.CreatedAt = time.Now()
}

func (ph *PasswordHistory) Create(ctx context.Context) {
	ph.FillDefaults()

	collection := helpers.GetCollection(config.UserDatabaseName, passwordHistoryCollectionName)
	_, err := collection.InsertOne(ctx, ph)

	if err != nil {
		log.Println(err)
	}
}

func (eh *EmailHistory) FillDefaults() {
	id := uuid.New().String()
	eh.ID = &id
	eh.CreatedAt = time.Now()
}

func (eh *EmailHistory) Create(ctx context.Context) {
	eh.FillDefaults()

	collection := helpers.GetCollection(config.UserDatabaseName, emailHistoryCollectionName)
	_, err := collection.InsertOne(ctx, eh)

	if err != nil {
		log.Println(err)
	}
}

func (pnh *PhoneNumberHistory) FillDefaults() {
	id := uuid.New().String()
	pnh.ID = &id
	pnh.CreatedAt = time.Now()
}

func (pnh *PhoneNumberHistory) Create(ctx context.Context) {
	pnh.FillDefaults()

	collection := helpers.GetCollection(config.UserDatabaseName, phoneNumberHistoryCollectionName)
	_, err := collection.InsertOne(ctx, pnh)

	if err != nil {
		log.Println(err)
	}
}
