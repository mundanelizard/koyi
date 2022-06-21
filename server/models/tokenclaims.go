package models

import (
	"context"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/services"
	"time"
)

type TokenClaims struct {
	ID           *string   `json:"id" bson:"id"`
	UserId       *string   `json:"userId" bson:"userId"`
	DeviceId     *string   `json:"deviceId" bson:"deviceId"`
	AccessToken  *string   `json:"accessToken" bson:"accessToken"`
	RefreshToken *string   `json:"refreshToken" bson:"refreshToken"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
}

func NewClaim(accessToken, refreshToken, deviceId, userId *string) *TokenClaims {
	return &TokenClaims{
		CreatedAt:    time.Now(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		DeviceId:     deviceId,
		UserId:       userId,
	}
}

func (tc *TokenClaims) Create(ctx context.Context) error {
	id := uuid.New().String()
	tc.ID = &id

	collection := services.GetCollection(config.UserDatabaseName, tokenClaimsCollectionName)
	_, err := collection.InsertOne(ctx, tc)

	return err
}
