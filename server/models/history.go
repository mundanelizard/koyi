package models

import (
	"context"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/services"
	"log"
	"time"
)

const (
	historyCollectionName = "history"
	passwordFieldName     = "password"
	emailFieldName        = "email"
	phoneNumberFieldName  = "phone-number"
)

type History struct {
	ID        string      `json:"id" bson:"id"`
	UserId    string      `json:"userId" bson:"userId"`
	Value     interface{} `json:"value" bson:"value"`
	Field     string      `json:"field" bson:"field"`
	Timestamp time.Time   `json:"createdAt" bson:"timestamp"`
}

func NewHistory(userId, field string, value interface{}) *History {
	id := uuid.New().String()
	return &History{
		ID:     id,
		UserId: userId,
		Value:  value,
		Field:  field,
	}
}

func (h *History) Create(ctx context.Context) {
	h.Timestamp = time.Now()
	collection := services.GetCollection(config.UserDatabaseName, historyCollectionName)
	_, err := collection.InsertOne(ctx, h)

	if err != nil {
		log.Println(err)
	}
}
