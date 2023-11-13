package models

import (
	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

func ConnectToDb() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("etm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")

	}
	return db
}
