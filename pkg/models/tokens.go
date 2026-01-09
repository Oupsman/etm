package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tokens struct {
	gorm.Model
	UUID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	User         Users
	UserID       uint   `json:"userid"`
	DeviceID     string `json:"deviceid"`
	Device       Devices
	RefreshToken string `json:"refresh_token"`
	Country      string `json:"country"`
}

func (db *DB) CreateToken(token *Tokens) error {
	result := db.Create(token)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *DB) UpdateToken(token *Tokens) error {
	result := db.Save(token)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *DB) DeleteToken(tokenID uint) error {
	result := db.Delete(&Tokens{}, tokenID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *DB) GetTokenByUUID(uuid uuid.UUID) (*Tokens, error) {
	var token Tokens
	result := db.Where("uuid = ?", uuid).First(&token)
	if result.Error != nil {
		return nil, result.Error
	}
	return &token, nil
}

func (db *DB) GetTokensByUserID(userID uint) ([]Tokens, error) {
	var tokens []Tokens
	result := db.Where("user_id = ?", userID).Find(&tokens)
	if result.Error != nil {
		return nil, result.Error
	}
	return tokens, nil
}
