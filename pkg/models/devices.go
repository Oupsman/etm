package models

import (
	"ETM/pkg/crypto"
	"ETM/pkg/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Devices struct {
	gorm.Model
	UUID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	User         Users
	UserID       uint                   `json:"userid"`
	DeviceID     crypto.EncryptedString `json:"deviceid"`
	DeviceName   crypto.EncryptedString `json:"devicename"`
	Trusted      bool                   `json:"trusted"`
	LastUsedAt   time.Time              `json:"lastusedat"`
	APIKeyHash   string                 `json:"-" gorm:"column:api_key_hash"`
	APIKeyPrefix string                 `json:"api_key_prefix" gorm:"column:api_key_prefix"`
}

func (d *Devices) BeforeCreate(_ *gorm.DB) error {
	if d.UUID == (uuid.UUID{}) {
		d.UUID = uuid.New()
	}
	return nil
}

func (db *DB) CreateDevice(device types.DeviceBody) error {
	var newDevice = Devices{
		DeviceID: crypto.EncryptedString(device.DeviceID),
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

func (db *DB) GetDeviceByID(deviceID uint) (*Devices, error) {
	var device Devices
	result := db.First(&device, deviceID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &device, nil
}

func (db *DB) GetDeviceByPrefix(prefix string) (*Devices, error) {
	var device Devices
	result := db.Preload("User").Where("api_key_prefix = ?", prefix).First(&device)
	if result.Error != nil {
		return nil, result.Error
	}
	return &device, nil
}

func (db *DB) UpdateDeviceLastUsed(deviceID uint) error {
	return db.Model(&Devices{}).Where("id = ?", deviceID).Update("last_used_at", time.Now()).Error
}

func (db *DB) RegisterDevice(userID uint, name, prefix, hash string) (*Devices, error) {
	device := Devices{
		UserID:       userID,
		DeviceName:   crypto.EncryptedString(name),
		APIKeyPrefix: prefix,
		APIKeyHash:   hash,
	}
	if err := db.Create(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}
