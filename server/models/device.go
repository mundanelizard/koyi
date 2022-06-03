package models

import (
	"context"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/helpers"
	"time"
)

const (
	deviceCollectionName = "devices"
)

type Device struct {
	ID        *string   `json:"id"`
	UserId    *string   `json:"userId"`
	Password  *string   `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

func (device *Device) Create(ctx *context.Context) error {
	collection := helpers.GetCollection(config.UserDatabaseName, deviceCollectionName)
	_, err := collection.InsertOne(*ctx, collection)
	return err
}
