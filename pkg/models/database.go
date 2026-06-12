package models

import (
	"database/sql"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	gorm.DB
}

var Db *gorm.DB

func ConnectToDB(driver string, dsn string) (*DB, error) {

	var err error
	var sqlDB *sql.DB
	var db *gorm.DB

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{PrepareStmt: true})

	if err != nil {
		return nil, err
	}

	sqlDB, err = db.DB()

	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(5)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	return &DB{*db}, nil
}

func (db *DB) Ping() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func CreateOrMigrate(db *DB) error {
	err := db.AutoMigrate(&Users{}, &Groups{}, &UserGroup{}, &Category{}, &CategoryGroup{}, &Tasks{}, &Keys{}, &Devices{}, &Tokens{})
	if err != nil {
		return err
	}
	return nil
}
