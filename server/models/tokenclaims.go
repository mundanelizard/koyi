package models

import (
	"context"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/services"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const (
	tokenClaimsCollectionName = "user-tokens"
)

type Tokens struct {
	ID           *string   `json:"id" bson:"id"`
	UserId       *string   `json:"userId" bson:"userId"`
	DeviceId     *string   `json:"deviceId" bson:"deviceId"`
	AccessToken  *string   `json:"accessToken" bson:"accessToken"`
	RefreshToken *string   `json:"refreshToken" bson:"refreshToken"`
	IsInvalid    *string   `json:"isInvalid" bson:"isInvalid"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
}

func NewTokens(deviceId string, user *User) (*Tokens, error) {
	at, err := NewClaim(accessToken, user)
	if err != nil {
		return nil, err
	}

	rt, err := NewClaim(refreshToken, user)
	if err != nil {
		return nil, err
	}

	tokens := &Tokens{
		CreatedAt:    time.Now(),
		AccessToken:  at,
		RefreshToken: rt,
		DeviceId:     &deviceId,
		UserId:       user.ID,
	}

	return tokens, nil
}

func FindClaim(ctx context.Context, userId, tokenType, token string) (*Tokens, error) {
	var fieldName = "accessToken"

	if tokenType == "refresh-token" {
		fieldName = "freshToken"
	}

	var tokens Tokens
	collection := services.GetCollection(config.UserDatabaseName, tokenClaimsCollectionName)
	err := collection.FindOne(ctx, bson.M{"userId": userId, fieldName: token}).Decode(&token)
	return &tokens, err
}

func (tc *Tokens) Create(ctx context.Context) error {
	id := uuid.New().String()
	tc.ID = &id

	collection := services.GetCollection(config.UserDatabaseName, tokenClaimsCollectionName)
	_, err := collection.InsertOne(ctx, tc)

	return err
}
