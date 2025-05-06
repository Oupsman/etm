package models

import (
	"ETM/pkg/vars"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func ConnectToDb() {
	var err error
	// Db, err := gorm.Open(sqlite.Open("etm.Db"), &gorm.Config{})
	dsn := "host=" + vars.DbHost + " user=" + vars.Username + " password=" + vars.Password + " dbname=" + vars.Database + " port=" + vars.DbPort + " sslmode=disable TimeZone=Europe/Paris"
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	sqlDB, err := Db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(5)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(10)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	if err != nil {
		panic("failed to connect database")
	}

}
