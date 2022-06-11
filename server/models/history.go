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
	phoneNumberHistoryCollectionName = "phone-number-history"
	emailHistoryCollectionName       = "email-history"
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

func NewPasswordHistory(userId, password *string) *PasswordHistory {
	id := uuid.New().String()
	return &PasswordHistory{
		UserId:   userId,
		Password: password,
		ID:       &id,
	}
}

func (ph *PasswordHistory) Create(ctx context.Context) {
	ph.CreatedAt = time.Now()
	collection := helpers.GetCollection(config.UserDatabaseName, passwordHistoryCollectionName)
	_, err := collection.InsertOne(ctx, ph)

	if err != nil {
		log.Println(err)
	}
}

func NewEmailHistory(userId, email *string) *EmailHistory {
	id := uuid.New().String()
	return &EmailHistory{
		UserId: userId,
		Email:  email,
		ID:     &id,
	}
}

func (eh *EmailHistory) Create(ctx context.Context) {
	eh.CreatedAt = time.Now()
	collection := helpers.GetCollection(config.UserDatabaseName, emailHistoryCollectionName)
	_, err := collection.InsertOne(ctx, eh)

	if err != nil {
		log.Println(err)
	}
}

func NewPhoneNumberHistory(userId *string, phoneNumber *PhoneNumber) *PhoneNumberHistory {
	id := uuid.New().String()
	return &PhoneNumberHistory{
		UserId:      userId,
		PhoneNumber: phoneNumber,
		ID:          &id,
	}
}

func (pnh *PhoneNumberHistory) Create(ctx context.Context) {
	pnh.CreatedAt = time.Now()
	collection := helpers.GetCollection(config.UserDatabaseName, phoneNumberHistoryCollectionName)
	_, err := collection.InsertOne(ctx, pnh)

	if err != nil {
		log.Println(err)
	}
}
