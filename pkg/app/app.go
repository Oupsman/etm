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
	err = models2.CreateOrMigrate(DB)
	if err != nil {
		return nil, err
	}

	var httpClient = &http.Client{}

	// var notif = notifications.New()

	return &App{DB: *DB, Client: httpClient, Logger: logger}, nil
}
