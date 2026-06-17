package models

import (
	"ETM/pkg/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	UUID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Username         string    `json:"username"`
	Password         string    `json:"-"`
	Gid              uint64    `json:"gid"`
	IsAdmin          string    `json:"isadmin"`
	Telegram         string    `json:"telegramconf"`
	Browser          string    `json:"browserconf"`
	Email            string    `json:"email"`
	OIDCSubject      string    `json:"oidc_subject,omitempty"  gorm:"column:oidc_subject;index"`
	OIDCProvider     string    `json:"oidc_provider,omitempty" gorm:"column:oidc_provider"`
	ActiveCategoryID uint      `json:"active_category_id"`
}

func (u *Users) BeforeCreate(_ *gorm.DB) error {
	if u.UUID == (uuid.UUID{}) {
		u.UUID = uuid.New()
	}
	return nil
}

func (db *DB) GetGroupUsers(c *gin.Context) ([]Users, error) {
	// var Db = ConnectToDb()
	var Users = []Users{}
	GroupID, _ := strconv.Atoi(c.Param("groupId"))
	result := db.Where("category_id = ?", GroupID).Find(&Users)
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

func (db *DB) CountUsers() (int64, error) {
	var count int64
	result := db.DB.Model(&Users{}).Count(&count)
	return count, result.Error
}

func (db *DB) CreateUser(user types.UserBody, isAdmin bool) (Users, error) {
	adminStr := ""
	if isAdmin {
		adminStr = "true"
	}
	newUser := Users{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		IsAdmin:  adminStr,
	}
	if err := db.DB.Create(&newUser).Error; err != nil {
		return Users{}, err
	}
	return newUser, nil
}

func (db *DB) SearchUsers(query string) ([]Users, error) {
	var users []Users
	like := "%" + query + "%"
	result := db.DB.Where("username ILIKE ? OR email ILIKE ?", like, like).Find(&users)
	return users, result.Error
}

// SeedAdmin promotes the oldest user to admin and creates an "Administrators"
// group for them if neither exists yet. Idempotent — safe to call on every startup.
func (db *DB) SeedAdmin() error {
	// Already have an admin — nothing to do.
	var count int64
	if err := db.DB.Model(&Users{}).Where("is_admin = ?", "true").Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// Find the first registered user (lowest primary key).
	var first Users
	if err := db.DB.Order("id ASC").First(&first).Error; err != nil {
		// No users yet — nothing to seed.
		return nil
	}

	// Promote to admin.
	if err := db.DB.Model(&first).Update("is_admin", "true").Error; err != nil {
		return err
	}

	// Create the Administrators group if it doesn't exist yet.
	var groupCount int64
	db.DB.Model(&Groups{}).Where("name = ? AND owner_id = ?", "Administrators", first.ID).Count(&groupCount)
	if groupCount == 0 {
		if _, err := db.CreateGroup("Administrators", first.ID); err != nil {
			return err
		}
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
