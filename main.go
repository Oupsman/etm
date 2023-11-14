package main

import (
	"ETM/models"
	"ETM/vars"
	"fmt"
	_ "gorm.io/driver/sqlite"
	"log"
	"net"
)

func initApp() {
	vars.Init()
	db := models.ConnectToDb()
	db.AutoMigrate(&models.Tasks{}, &models.Category{})

	// Check if a category exists
	var category = models.Category{}
	result := db.First(&category)
	if result.Error != nil {
		// Create a category
		category.Name = "Default"
		db.Create(&category)
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
