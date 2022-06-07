package helpers

import (
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/models"
	"net/http"
)

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

func ExtractDevice(r *http.Request, userId *string) *models.Device {
	// todo => implement this
	id := uuid.New().String()
	return &models.Device{
		ID:     &id,
		UserId: userId,
	}
}
