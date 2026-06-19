package models

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (db *DB) GetUserByOIDCSubject(subject string) (Users, error) {
	var user Users
	result := db.DB.Where("oidc_subject = ?", subject).First(&user)
	return user, result.Error
}

func (db *DB) LinkOIDCToUser(userID uint, subject, provider string) error {
	return db.DB.Model(&Users{}).Where("id = ?", userID).
		Updates(map[string]interface{}{"oidc_subject": subject, "oidc_provider": provider}).Error
}

func (db *DB) UnlinkOIDC(userID uint) error {
	return db.DB.Model(&Users{}).Where("id = ?", userID).
		Updates(map[string]interface{}{"oidc_subject": "", "oidc_provider": ""}).Error
}

// FindOrCreateOIDCUser looks up a user by their OIDC subject claim,
// creating them if they don't exist yet. The second return value is true
// when the user was just created (false for returning users).
func (db *DB) FindOrCreateOIDCUser(subject, email, username string) (Users, bool, error) {
	var user Users

	result := db.Where("oidc_subject = ?", subject).First(&user)
	if result.Error == nil {
		return user, false, nil
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return Users{}, false, result.Error
	}

	// First login: create the user (no password)
	user = Users{
		UUID:         uuid.New(),
		Username:     username,
		Email:        email,
		OIDCSubject:  subject,
		OIDCProvider: "oidc",
	}
	if err := db.Create(&user).Error; err != nil {
		return Users{}, false, err
	}
	return user, true, nil
}
