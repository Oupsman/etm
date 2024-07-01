package models

import "gorm.io/gorm"

type Keys struct {
	gorm.Model
	Pubkey  string
	Privkey string
}

func SaveKeys(keys *Keys) error {
	result := Db.Create(&keys)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetKeys() (*Keys, error) {
	var keys = Keys{}
	result := Db.First(&keys)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &keys, nil
}
