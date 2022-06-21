package testutils

import (
	"context"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/services"
	"log"
)

func ClearDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), config.AverageServerTimeout)
	defer cancel()

	err := services.GetDatabase(config.UserDatabaseName).Drop(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}
