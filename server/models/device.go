package models

import (
	"context"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/helpers"
	"net/http"
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

/**
todo => extract and store
tenant_id: the course creator school the visitor is checking
raw: the raw ua
type: desktop / mobile / tablet / bot / other
browser_name
browser_version
os_name: Android / IOS / Windows / Mac
os_version: OS Version
hardware_details: hstore containing memory, processor, device_model, device_name
connection_speed: hstore containing downlink_max, connection_type
*/

func ExtractDevice(r *http.Request, userId *string) *Device {
	// todo => implement this
	id := uuid.New().String()
	return &Device{
		ID:     &id,
		UserId: userId,
	}
}
