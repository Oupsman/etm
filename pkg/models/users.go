package models

import (
	"ETM/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Users struct {
	gorm.Model
	UUID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name     string    `json:"username"`
	Password string    `json:"password"`
	Gid      uint64    `json:"gid"`
	IsAdmin  string    `json:"isadmin"`
	Telegram string    `json:"telegramconf"`
	Browser  string    `json:"browserconf"`
	Email    string    `json:"email"`
}

type Groups struct {
	gorm.Model
	Name  string `json:"name"`
	Users []Users
}

func (db *DB) GetGroupUsers(c *gin.Context) ([]Users, error) {
	// var Db = ConnectToDb()
	var Users = []Users{}
	GroupID, _ := strconv.Atoi(c.Param("groupId"))
	result := Db.Where("category_id = ?", GroupID).Find(&Users)
	if result.Error != nil {
		return nil, result.Error
	}
	return Users, nil
}

func (db *DB) GetAllUsers() ([]Users, error) {
	var Users = []Users{}
	result := db.Find(&Users)
	if result.Error != nil {
		return nil, result.Error
	}
	return Users, nil
}

func (db *DB) GetUserDetails(c *gin.Context) {
	var User = Users{}
	UserID, _ := strconv.Atoi(c.Param("userId"))
	result := db.Where("user_id = ?", UserID).Find(&User)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unable to list tasks"})
		return
	}
	c.JSON(http.StatusOK, User)
}

func (db *DB) GetUser(userID uint) (Users, error) {
	var user Users
	result := db.First(&user, userID)
	if result.Error != nil {
		return Users{}, result.Error
	}
	return user, nil
}

func (db *DB) CreateUser(user types.UserBody) error {
	var newUser = Users{
		Name:     user.Username,
		Password: user.Password,
		Email:    user.Email,
	}
	result := db.Create(&newUser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *DB) UpdateUser(user Users) error {

	result := db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *DB) GetUserByUUID(uuid uuid.UUID) (Users, error) {
	var user Users
	result := db.Where("uuid = ?", uuid).First(&user)
	if result.Error != nil {
		return Users{}, result.Error
	}
	return user, nil
}
