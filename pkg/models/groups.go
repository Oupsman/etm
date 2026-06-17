package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Groups struct {
	gorm.Model
	Name    string `json:"name"`
	OwnerID uint   `json:"owner_id"`
}

// UserGroup is the join table between Users and Groups.
// Role is one of: "owner", "writer", "reader".
type UserGroup struct {
	UserID  uint   `json:"user_id" gorm:"primaryKey"`
	GroupID uint   `json:"group_id" gorm:"primaryKey"`
	Role    string `json:"role"`
}

type GroupMemberDetail struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func (db *DB) CreateGroup(name string, ownerID uint) (Groups, error) {
	var group Groups
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		group = Groups{Name: name, OwnerID: ownerID}
		if err := tx.Create(&group).Error; err != nil {
			return err
		}
		return tx.Create(&UserGroup{UserID: ownerID, GroupID: group.ID, Role: "owner"}).Error
	})
	return group, err
}

func (db *DB) GetGroupsByUser(userID uint) ([]Groups, error) {
	var groups []Groups
	result := db.DB.Table("groups").
		Joins("JOIN user_groups ON user_groups.group_id = groups.id").
		Where("user_groups.user_id = ?", userID).
		Where("groups.deleted_at IS NULL").
		Find(&groups)
	return groups, result.Error
}

func (db *DB) GetGroup(groupID uint) (Groups, error) {
	var group Groups
	result := db.DB.First(&group, groupID)
	return group, result.Error
}

func (db *DB) DeleteGroup(groupID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("group_id = ?", groupID).Delete(&UserGroup{}).Error; err != nil {
			return err
		}
		if err := tx.Where("group_id = ?", groupID).Delete(&CategoryGroup{}).Error; err != nil {
			return err
		}
		return tx.Delete(&Groups{}, groupID).Error
	})
}

func (db *DB) IsGroupOwner(groupID, userID uint) bool {
	var group Groups
	if err := db.DB.First(&group, groupID).Error; err != nil {
		return false
	}
	if group.OwnerID == userID {
		return true
	}
	// Also accept members explicitly granted the "owner" role.
	var count int64
	db.DB.Table("user_groups").
		Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, "owner").
		Count(&count)
	return count > 0
}

func (db *DB) IsGroupMember(groupID, userID uint) bool {
	var count int64
	db.DB.Table("user_groups").
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Count(&count)
	return count > 0
}

func (db *DB) AddGroupMember(groupID, userID uint, role string) error {
	if role != "writer" && role != "reader" && role != "owner" {
		return fmt.Errorf("invalid role: must be 'owner', 'writer' or 'reader'")
	}
	var count int64
	db.DB.Table("user_groups").Where("group_id = ? AND user_id = ?", groupID, userID).Count(&count)
	if count > 0 {
		return fmt.Errorf("user is already a member of this group")
	}
	return db.DB.Create(&UserGroup{UserID: userID, GroupID: groupID, Role: role}).Error
}

func (db *DB) UpdateGroupMember(groupID, userID uint, role string) error {
	if role != "writer" && role != "reader" && role != "owner" {
		return fmt.Errorf("invalid role: must be 'owner', 'writer' or 'reader'")
	}
	return db.DB.Model(&UserGroup{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Update("role", role).Error
}

func (db *DB) RemoveGroupMember(groupID, userID uint) error {
	return db.DB.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&UserGroup{}).Error
}

func (db *DB) GetAllGroups() ([]Groups, error) {
	var groups []Groups
	result := db.DB.Where("deleted_at IS NULL").Find(&groups)
	return groups, result.Error
}

func (db *DB) GetGroupMembers(groupID uint) ([]GroupMemberDetail, error) {
	var members []GroupMemberDetail
	result := db.DB.Table("user_groups").
		Select("user_groups.user_id, users.username, users.email, user_groups.role").
		Joins("JOIN users ON users.id = user_groups.user_id").
		Where("user_groups.group_id = ?", groupID).
		Where("users.deleted_at IS NULL").
		Scan(&members)
	return members, result.Error
}