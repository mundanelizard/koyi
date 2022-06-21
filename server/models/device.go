package models

import (
	"context"
	"errors"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

const (
	deviceCollectionName = "devices"
)

type Device struct {
	ID      string `json:"id" bson:"id"`           // 000-0000-000-000-000 | build-number
	UserId  string `json:"userId" bson:"userId"`   // 000-0000-000-000-000
	Name    string `json:"name" bson:"name"`       // Samsung Galaxy Note 10
	OS      string `json:"os" bson:"os"`           // Android | IOS | Windows | Mac
	Version string `json:"version" bson:"version"` // Android 2.41
	Type    string `json:"type"`                   // web - desktop - mobile - bot - other
}

func (device *Device) Create(ctx context.Context) error {
	collection := helpers.GetCollection(config.UserDatabaseName, deviceCollectionName)
	_, err := collection.InsertOne(ctx, collection)
	return err
}

// Exists checks if a user exists in the database.
func (device *Device) Exists(ctx context.Context) (bool, error) {
	if device.ID == "" {
		return false, errors.New("invalid device id")
	}

	var count int64
	var err error

	count, err = CountDevice(ctx, map[string]string{"id": device.ID})

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

func ExtractDevice(r *http.Request) *Device {
	return &Device{
		ID:      r.Header.Get("device-id"),
		OS:      r.Header.Get("platform"),
		Type:    r.Header.Get("platform-type"),
		Version: r.Header.Get("platform-version"),
		Name:    r.Header.Get("platform-name"),
	}
}

func ExtractAndCreateDevice(ctx context.Context, r *http.Request, userId string) *Device {
	device := ExtractDevice(r)
	device.UserId = userId
	err := device.Create(ctx)

	log.Println(err)

	return device
}

func FindDevice(ctx context.Context, deviceId string) (*Device, error) {
	var device Device

	collection := helpers.GetCollection(config.UserDatabaseName, intentsCollectionName)
	err := collection.FindOne(ctx, bson.M{"id": deviceId}).Decode(&device)

	return &device, err
}
