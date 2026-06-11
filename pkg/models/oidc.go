package models

import "gorm.io/gorm"

// Call this from CreateOrMigrate in database.go alongside existing AutoMigrate calls
func MigrateOIDC(db *gorm.DB) error {
	return db.AutoMigrate(&Users{})
}

// FindOrCreateOIDCUser looks up a user by their OIDC subject claim,
// creating them if they don't exist yet.
func (db *DB) FindOrCreateOIDCUser(subject, email, username string) (Users, error) {
	var user Users

	result := db.Where("oidc_subject = ?", subject).First(&user)
	if result.Error == nil {
		// Existing OIDC user
		return user, nil
	}
	if result.Error != gorm.ErrRecordNotFound {
		return Users{}, result.Error
	}

	// First login: create the user (no password)
	user = Users{
		Username:     username,
		Email:        email,
		OIDCSubject:  subject,
		OIDCProvider: "oidc",
	}
	if err := db.Create(&user).Error; err != nil {
		return Users{}, err
	}
	return user, nil
}
