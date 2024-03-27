package models

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Name     string `json:"name"`
	Password string `json:"password"`
	Gid      uint   `json:"gid"`
	IsAdmin  string `json:"isadmin"`
}

type Groups struct {
	gorm.Model
	Name  string `json:"name"`
	Users []Users
}

func GetGroupUsers(c *gin.Context) {
	// var Db = ConnectToDb()
	var Users = []Users{}
	GroupID, _ := strconv.Atoi(c.Param("groupId"))
	result := Db.Where("category_id = ?", GroupID).Find(&Users)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unable to list tasks"})
		return
	}
	c.JSON(http.StatusOK, Users)
}

func GetAllUsers(c *gin.Context) {
	// var Db = ConnectToDb()
	var Users = []Users{}
	result := Db.Find(&Users)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unable to list tasks"})
		return
	}
	c.JSON(http.StatusOK, Users)
}

func GetUserDetails(c *gin.Context) {
	var User = Users{}
	UserID, _ := strconv.Atoi(c.Param("userId"))
	result := Db.Where("user_id = ?", UserID).Find(&User)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unable to list tasks"})
		return
	}
	c.JSON(http.StatusOK, User)

}
