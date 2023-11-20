package models

import (
	"ETM/vars"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDb() *gorm.DB {
	// db, err := gorm.Open(sqlite.Open("etm.db"), &gorm.Config{})
	dsn := "host=" + vars.DbHost + " user=" + vars.Username + " password=" + vars.Password + " dbname=" + vars.Database + " port=" + vars.DbPort + " sslmode=disable TimeZone=Europe/Paris"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
