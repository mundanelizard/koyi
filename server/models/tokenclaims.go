package models

import (
	"context"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/helpers"
	"time"
)

type TokenClaims struct {
	ID           *string    `json:"id" bson:"id"`
	DeviceId     *string    `json:"deviceId" bson:"deviceId"`
	AccessToken  *string    `json:"accessToken" bson:"accessToken"`
	RefreshToken *string    `json:"refreshToken" bson:"refreshToken"`
	RefreshClaim *UserClaim `json:"refreshClaim" bson:"refreshClaim"`
	AccessClaim  *UserClaim `json:"accessClaim" bson:"accessClaim"`
	CreatedAt    time.Time  `json:"createdAt" bson:"createdAt"`
}

func NewTokenClaim(accessToken, refreshToken *string, refreshClaim, accessClaim *UserClaim, deviceId *string) *TokenClaims {
	return &TokenClaims{
		CreatedAt:    time.Now(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,

		// todo => I may remove thing because i think it's redundant.
		RefreshClaim: refreshClaim,
		AccessClaim:  accessClaim,

		DeviceId: deviceId,
	}
}

func (tc *TokenClaims) Create(ctx context.Context) error {
	id := uuid.New().String()
	tc.ID = &id

	collection := helpers.GetCollection(config.UserDatabaseName, tokenClaimsCollectionName)
	_, err := collection.InsertOne(ctx, tc)

	return err
}
