package main

import (
	"ETM/pkg/app"
	"ETM/pkg/controllers"
	"ETM/pkg/models"
	"ETM/pkg/vars"
	"ETM/pkg/webserver"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_ "gorm.io/driver/postgres"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	vars.Init()
	driver := "postgres"
	dsn := vars.Dsn
	App, err := app.NewApp(log.Logger, driver, dsn)

	if err != nil {
		panic(err)
	}
	// Generate VAPID keys on first run if none exist yet.
	if keys, _ := models.GetKeys(); keys == nil {
		if err := controllers.GenerateVapidKeys(); err != nil {
			log.Warn().Err(err).Msg("failed to generate VAPID keys — browser notifications will not work")
		} else {
			log.Info().Msg("VAPID keys generated")
		}
	}

	// Background goroutine: send push notifications for tasks due within 1 hour.
	// Checks every hour; deduplicates within a run to avoid spamming.
	go func() {
		notified := make(map[uint]struct{})
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			tasks, err := App.DB.GetTasksDueSoon(24 * time.Hour)
			if err != nil {
				log.Error().Err(err).Msg("notification worker: query failed")
				continue
			}
			// Remove tasks no longer in the due-soon window from the dedup set.
			inWindow := make(map[uint]struct{}, len(tasks))
			for _, t := range tasks {
				inWindow[t.ID] = struct{}{}
			}
			for id := range notified {
				if _, still := inWindow[id]; !still {
					delete(notified, id)
				}
			}
			for _, task := range tasks {
				if _, already := notified[task.ID]; already {
					continue
				}
				if task.User.Browser == "" {
					continue
				}
				msg, _ := json.Marshal(map[string]string{
					"title": "Task due soon",
					"body":  task.Name + " is due within 24 hours",
				})
				if err := controllers.BrowserSend(string(msg), task.User.Browser); err != nil {
					log.Warn().Err(err).Str("task", task.Name).Msg("notification worker: push failed")
				} else {
					notified[task.ID] = struct{}{}
				}
			}
		}
	}()

	fmt.Println("Starting Eisenhower TaskCard Manager")
	if err := controllers.InitOIDC(); err != nil {
		log.Fatal().Err(err).Msg("failed to initialize OIDC")
	}
	fmt.Printf("Listening on %s:%s\n", vars.Host, vars.Port)

	addr := net.JoinHostPort(vars.Host, vars.Port)
	if err := webserver.RunHttp(addr, App); err != nil {
		log.Fatal()
	}

}
