package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/models"
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

func ExtractDeviceDetailsFromContext(c *gin.Context) *models.Device {
	// todo => implement this
	id := uuid.New().String()
	return &models.Device{
		ID: &id,
	}
}
