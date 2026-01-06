package models

import (
	"ETM/pkg/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Devices struct {
	gorm.Model
	UUID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	User       Users
	UserID     uint      `json:"userid"`
	DeviceID   string    `json:"deviceid"`
	LastUsedAt time.Time `json:"lastusedat"`
}

func (db *DB) CreateDevice(device types.DeviceBody) error {
	var newDevice = Devices{
		DeviceID: device.DeviceID,
		UserID:   device.UserID,
	}
	result := db.Create(&newDevice)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *DB) GetDevices(userID uint) ([]Devices, error) {
	var devices []Devices
	result := db.Where("user_id = ?", userID).Find(&devices)
	if result.Error != nil {
		return nil, result.Error
	}
	return devices, nil
}

func (db *DB) DeleteDevice(deviceID uint) error {
	result := db.Delete(&Devices{}, deviceID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
