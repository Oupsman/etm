package models

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
)

type Category struct {
	gorm.Model
	Name  string `json:"name" binding:"required"`
	Color string `json:"color" binding:"required"`
}

type CategoryBody struct {
	Name  string
	Color string
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

	var category = Category{}

	err := c.ShouldBindJSON(&category)
	switch {
	case errors.Is(err, io.EOF):
		fmt.Println("Error :", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing request body"})
		return
	case err != nil:
		fmt.Println("Error :", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Create(&category)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}
