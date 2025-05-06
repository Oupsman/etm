package main

import (
	"ETM/pkg/app"
	"ETM/pkg/vars"
	"ETM/pkg/webserver"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_ "gorm.io/driver/postgres"
	"net"
	"os"
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
	//err = controllers.GenerateVapidKeys()

	fmt.Println("Starting Eisenhower TaskCard Manager")
	// 	fmt.Printf("Username: %s\n", vars.Username)
	//	fmt.Printf("Token: %s\n", vars.Token)
	// 	fmt.Printf("Connection String: %s\n", vars.ConnectionString)
	/*
		go func() {
			for {
				// Check every 12 hours if we have tasks with due date expired
				tasks, err := models2.GetActiveTasks()
				if err != nil {
					log.Fatal(err)
				}
				for _, task := range tasks {
					now := time.Now()
					if task.DueDate.Sub(now) < 0 {
						message := "TaskCard " + task.Name + " is overdue"
						user, err := models2.GetUser(task.UserID)
						if err != nil {
							log.Fatal(err)
						}
						err = controllers.BrowserSend(message, user.Browser)
						if err != nil {
							log.Fatal(err)
						}
					}
					time.Sleep(2 * time.Second)
				}
				time.Sleep(12 * time.Hour)
			}
		}()*/

	fmt.Printf("Listening on %s:%s\n", vars.Host, vars.Port)

	addr := net.JoinHostPort(vars.Host, vars.Port)
	if err := webserver.RunHttp(addr, App); err != nil {
		log.Fatal()
	}

}
