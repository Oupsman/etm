package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Category struct {
	gorm.Model
	Name string
}

func GetCategories(c *gin.Context) {
	var db = ConnectToDb()
	var Categories = []Category{}
	result := db.Find(&Categories)

	if result.Error != nil {
		c.JSON(http.StatusForbidden, Categories)
	}
	c.JSON(http.StatusOK, Categories)
}

func CreateCategory(c *gin.Context) {
	var db = ConnectToDb()
	var name = c.Query("name")
	var category = Category{
		Name: name,
	}

	result := db.Create(&category)
	if result.Error != nil {
		c.JSON(http.StatusForbidden, category)
	}

	c.JSON(http.StatusOK, category)
}
