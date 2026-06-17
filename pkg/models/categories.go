package models

import (
	"ETM/pkg/types"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name   string `json:"name" binding:"required"`
	Color  string `json:"color" binding:"required"`
	UserID uint   `json:"userid" binding:"required"`
	Active bool   `json:"active" binding:"required"`
	User   Users
}

// CategoryGroup is the join table between Category and Groups.
type CategoryGroup struct {
	CategoryID uint `json:"category_id" gorm:"primaryKey"`
	GroupID    uint `json:"group_id" gorm:"primaryKey"`
}

// CategoryWithSharing wraps a Category with sharing metadata for API responses.
type CategoryWithSharing struct {
	Category
	Shared    bool   `json:"shared"`
	GroupName string `json:"group_name,omitempty"`
	GroupID   uint   `json:"group_id,omitempty"`
	UserRole  string `json:"user_role,omitempty"` // viewer / editor / owner; empty when the user owns the category
}

func (db *DB) GetCategory(categoryID uint) (Category, error) {
	var cat Category
	result := db.DB.First(&cat, categoryID)
	return cat, result.Error
}

func (db *DB) GetCategories(userID uint) ([]CategoryWithSharing, error) {
	result := []CategoryWithSharing{}

	var owned []Category
	if err := db.DB.Where("user_id = ?", userID).Find(&owned).Error; err != nil {
		return nil, err
	}
	for _, c := range owned {
		result = append(result, CategoryWithSharing{Category: c})
	}

	var userGroups []UserGroup
	if err := db.DB.Where("user_id = ?", userID).Find(&userGroups).Error; err != nil {
		return nil, err
	}
	if len(userGroups) == 0 {
		return result, nil
	}

	groupIDs := make([]uint, len(userGroups))
	for i, ug := range userGroups {
		groupIDs[i] = ug.GroupID
	}

	type sharedRow struct {
		ID        uint
		GroupID   uint
		GroupName string
		UserRole  string
	}
	var sharedRows []sharedRow
	if err := db.DB.Table("category_groups").
		Select("category_groups.category_id as id, category_groups.group_id, groups.name as group_name, user_groups.role as user_role").
		Joins("JOIN groups ON groups.id = category_groups.group_id").
		Joins("JOIN user_groups ON user_groups.group_id = category_groups.group_id AND user_groups.user_id = ?", userID).
		Where("category_groups.group_id IN ?", groupIDs).
		Where("groups.deleted_at IS NULL").
		Scan(&sharedRows).Error; err != nil {
		return nil, err
	}

	if len(sharedRows) == 0 {
		return result, nil
	}

	rolePriority := map[string]int{"reader": 1, "writer": 2, "owner": 3}
	rowByID := make(map[uint]sharedRow, len(sharedRows))
	sharedCatIDs := make([]uint, 0, len(sharedRows))
	for _, row := range sharedRows {
		existing, seen := rowByID[row.ID]
		if !seen {
			sharedCatIDs = append(sharedCatIDs, row.ID)
			rowByID[row.ID] = row
		} else if rolePriority[row.UserRole] > rolePriority[existing.UserRole] {
			// Keep the highest-privilege role when shared via multiple groups.
			rowByID[row.ID] = row
		}
	}

	var sharedCats []Category
	if err := db.DB.Where("id IN ? AND user_id != ?", sharedCatIDs, userID).Find(&sharedCats).Error; err != nil {
		return nil, err
	}
	for _, cat := range sharedCats {
		row := rowByID[cat.ID]
		result = append(result, CategoryWithSharing{
			Category:  cat,
			Shared:    true,
			GroupName: row.GroupName,
			GroupID:   row.GroupID,
			UserRole:  row.UserRole,
		})
	}

	return result, nil
}

// CanUserViewCategory returns true if userID owns the category or is a member
// of any group it has been shared with.
func (db *DB) CanUserViewCategory(userID, categoryID uint) bool {
	var cat Category
	if err := db.DB.First(&cat, categoryID).Error; err != nil {
		return false
	}
	if cat.UserID == userID {
		return true
	}
	var count int64
	db.DB.Table("category_groups").
		Joins("JOIN user_groups ON user_groups.group_id = category_groups.group_id").
		Where("category_groups.category_id = ? AND user_groups.user_id = ?", categoryID, userID).
		Count(&count)
	return count > 0
}

// CanUserEditCategory returns true if userID owns the category or has a
// "writer" or "owner" role in a group the category is shared with.
func (db *DB) CanUserEditCategory(userID, categoryID uint) bool {
	var cat Category
	if err := db.DB.First(&cat, categoryID).Error; err != nil {
		return false
	}
	if cat.UserID == userID {
		return true
	}
	var count int64
	db.DB.Table("category_groups").
		Joins("JOIN user_groups ON user_groups.group_id = category_groups.group_id").
		Where("category_groups.category_id = ? AND user_groups.user_id = ? AND user_groups.role IN ?",
			categoryID, userID, []string{"writer", "owner"}).
		Count(&count)
	return count > 0
}

func (db *DB) CreateCategory(newCategory types.CategoryBody, user Users) (Category, error) {
	category := Category{
		Name:   newCategory.Name,
		Color:  newCategory.Color,
		UserID: user.ID,
		Active: newCategory.Active,
	}
	result := db.DB.Create(&category)
	return category, result.Error
}

func (db *DB) UpdateCategory(category Category) error {
	return db.DB.Save(&category).Error
}

func (db *DB) ShareCategoryWithGroup(categoryID, groupID uint) error {
	return db.DB.Create(&CategoryGroup{CategoryID: categoryID, GroupID: groupID}).Error
}

func (db *DB) UnshareCategoryFromGroup(categoryID, groupID uint) error {
	return db.DB.Where("category_id = ? AND group_id = ?", categoryID, groupID).
		Delete(&CategoryGroup{}).Error
}

func (db *DB) GetCategoryShares(categoryID uint) ([]Groups, error) {
	var groups []Groups
	result := db.DB.Table("groups").
		Joins("JOIN category_groups ON category_groups.group_id = groups.id").
		Where("category_groups.category_id = ?", categoryID).
		Where("groups.deleted_at IS NULL").
		Find(&groups)
	return groups, result.Error
}