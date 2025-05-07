package models

// Package model for Categories
import (
	"ETM/pkg/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name   string `json:"name" binding:"required"`
	Color  string `json:"color" binding:"required"`
	UserID uint   `json:"userid" binding:"required"`
	Active bool   `json:"active" binding:"required"`
	User   Users
}

func (db *DB) GetCategories(UserUUID uuid.UUID) ([]Category, error) {
	// var Db = ConnectToDb()
	var Categories = []Category{}
	result := db.Find(&Categories).Where("useruuid = ?", UserUUID)

	if result.Error != nil {
		return nil, result.Error
	}
	return Categories, nil
}

func (db *DB) CreateCategory(newCategory types.CategoryBody, user Users) (Category, error) {
	//  var Db = ConnectToDb()
	var category = Category{
		Name:   newCategory.Name,
		Color:  newCategory.Color,
		UserID: user.ID,
		Active: newCategory.Active,
	}
	result := db.Create(&category)
	if result.Error != nil {
		return Category{}, result.Error
	}

	return category, nil
}

func (db *DB) UpdateCategory(category Category) error {
	// var Db = ConnectToDb()
	result := db.Save(&category)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
