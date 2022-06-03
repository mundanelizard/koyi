package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/models"
)

func ExtractDeviceDetailsFromContext(c *gin.Context) *models.Device {
	// todo => implement this
	id := uuid.New().String()
	return &models.Device{
		ID: &id,
	}
}
