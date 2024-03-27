package main

import (
	"ETM/models"
	"ETM/vars"
	"fmt"
	"log"
	"net"

	_ "gorm.io/driver/postgres"
)

func initApp() {
	vars.Init()
	models.ConnectToDb()
	models.Db.AutoMigrate(&models.Tasks{}, &models.Category{}, &models.Users{})

	// Check if a category exists
	var category = models.Category{}
	result := models.Db.First(&category)
	if result.Error != nil {
		// Create a category
		category.Name = "Default"
		models.Db.Create(&category)
	}

}

func main() {
	initApp()

	fmt.Println("Starting Eisenhower Task Manager")
	// 	fmt.Printf("Username: %s\n", vars.Username)
	//	fmt.Printf("Token: %s\n", vars.Token)
	// 	fmt.Printf("Connection String: %s\n", vars.ConnectionString)
	fmt.Printf("Listening on %s:%s\n", vars.Host, vars.Port)

	addr := net.JoinHostPort(vars.Host, vars.Port)
	if err := runHttp(addr); err != nil {
		log.Fatal(err)
	}

}
