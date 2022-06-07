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

func (device *Device) Create(ctx context.Context) error {
	collection := helpers.GetCollection(config.UserDatabaseName, deviceCollectionName)
	_, err := collection.InsertOne(ctx, collection)
	return err
}

// Exists checks if a user exists in the database.
func (device *Device) Exists(ctx context.Context) (bool, error) {
	var count int64
	var err error

	count, err = CountDevice(ctx, map[string]string{"id": *device.ID})

	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func CountDevice(ctx context.Context, filter interface{}) (int64, error) {
	collection := helpers.GetCollection(config.UserDatabaseName, deviceCollectionName)
	count, err := collection.CountDocuments(ctx, filter)
	return count, err
}
