package models

import (
	"context"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/helpers"
	"time"
)

type TokenClaims struct {
	ID           *string    `json:"id"`
	DeviceId     *string    `json:"deviceId"`
	AccessToken  *string    `json:"accessToken"`
	RefreshToken *string    `json:"refreshToken"`
	RefreshClaim *UserClaim `json:"refreshClaim"`
	AccessClaim  *UserClaim `json:"accessClaim"`
	CreatedAt    time.Time  `json:"createdAt"`
}

func (tc *TokenClaims) Create(ctx context.Context) error {
	id := uuid.New().String()
	tc.ID = &id

	collection := helpers.GetCollection(config.UserDatabaseName, tokenClaimsCollectionName)
	_, err := collection.InsertOne(ctx, tc)

	return err
}
