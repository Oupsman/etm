package models

import (
	"ETM/types"
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
	Telegram string `json:"telegramconf"`
	Browser  string `json:"browserconf"`
	Email    string `json:"email"`
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

func GetUser(userID uint) (Users, error) {
	var user Users
	result := Db.First(&user, userID)
	if result.Error != nil {
		return Users{}, result.Error
	}
	return user, nil
}

func CreateUser(user types.UserBody) error {
	var newUser = Users{
		Name:     user.Username,
		Password: user.Password,
		Email:    user.Email,
	}
	result := Db.Create(&newUser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateUser(user Users) error {

	result := Db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
