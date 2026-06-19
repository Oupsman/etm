package app

import (
	models2 "ETM/pkg/models"
	"github.com/rs/zerolog"
	"net/http"
)

type App struct {
	DB     models2.DB
	Client *http.Client
	Logger zerolog.Logger
	//	Notifications *notifications.Notifs
}

func NewApp(logger zerolog.Logger, driver string, dsn string) (*App, error) {
	DB, err := models2.ConnectToDB(driver, dsn)
	if err != nil {
		return nil, err
	}
	if err = models2.CreateOrMigrate(DB); err != nil {
		return nil, err
	}
	if err = DB.SeedAdmin(); err != nil {
		logger.Warn().Err(err).Msg("admin seed failed")
	}
	if err = DB.SeedDemoContentForAllUsers(); err != nil {
		logger.Warn().Err(err).Msg("demo content seed failed")
	}

	var httpClient = &http.Client{}

	return &App{DB: *DB, Client: httpClient, Logger: logger}, nil
}
