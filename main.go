package main

import (
	"ETM/controllers"
	"ETM/models"
	"ETM/vars"
	"fmt"
	_ "gorm.io/driver/postgres"
	"log"
	"net"
	"time"
)

func initApp() error {
	vars.Init()
	models.ConnectToDb()
	err := models.Db.AutoMigrate(&models.Tasks{}, &models.Category{}, &models.Users{}, &models.Keys{})
	if err != nil {
		return err
	}
	// Check if a category exists
	var category = models.Category{}
	result := models.Db.First(&category)
	if result.Error != nil {
		// Create a category
		category.Name = "Default"
		models.Db.Create(&category)
	}
	return nil
}

func main() {
	err := initApp()
	if err != nil {
		panic(err)
	}
	err = controllers.GenerateVapidKeys()
	fmt.Println("Starting Eisenhower Task Manager")
	// 	fmt.Printf("Username: %s\n", vars.Username)
	//	fmt.Printf("Token: %s\n", vars.Token)
	// 	fmt.Printf("Connection String: %s\n", vars.ConnectionString)

	go func() {
		for {
			// Check every 12 hours if we have tasks with due date expired
			tasks, err := models.GetActiveTasks()
			if err != nil {
				log.Fatal(err)
			}
			for _, task := range tasks {
				now := time.Now()
				if task.DueDate.Sub(now) < 0 {
					message := "Task " + task.Name + " is overdue"
					user, err := models.GetUser(task.UserID)
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
	}()

	fmt.Printf("Listening on %s:%s\n", vars.Host, vars.Port)

	addr := net.JoinHostPort(vars.Host, vars.Port)
	if err := runHttp(addr); err != nil {
		log.Fatal(err)
	}

}
