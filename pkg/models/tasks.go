package models

import (
	"ETM/pkg/crypto"
	"time"

	"gorm.io/gorm"
)

type Tasks struct {
	gorm.Model
	Name        crypto.EncryptedString `json:"name"`
	Comment     crypto.EncryptedString `json:"comment"`
	Link        crypto.EncryptedString `json:"link"`
	IsCompleted bool                   `json:"iscompleted"`
	IsStarted   bool                   `json:"isstarted"`
	IsBackLog   bool                   `json:"isbacklog"`
	CategoryID  uint                   `json:"category-id"`
	Category    Category
	Priority    bool      `json:"priority"`
	Urgency     bool      `json:"urgency"`
	DueDate     time.Time `json:"duedate"`
	User        Users
	UserID      uint `json:"userid"`
}

func (db *DB) GetTasks(UserID uint, CategoryId int) ([]Tasks, error) {
	// var Db = ConnectToDb()
	var Tasks = []Tasks{}
	result := db.Where("user_id = ? AND category_id = ?", UserID, CategoryId).Find(&Tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return Tasks, nil
}

func (db *DB) GetTask(TaskID int) (*Tasks, error) {
	// var Db = ConnectToDb()
	var task = Tasks{}
	result := db.First(&task, TaskID)

	if result.Error != nil {

		return nil, result.Error
	}

	return &task, nil
}

func (db *DB) CreateTask(task *Tasks) error {

	result := db.Create(&task)
	if result.Error != nil {

		return result.Error
	}

	return nil
}

func (db *DB) UpdateTask(task *Tasks) error {
	// var Db = ConnectToDb()

	result := db.Save(&task)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) GetActiveTasks() ([]Tasks, error) {

	var tasks []Tasks

	result := Db.Model(&Tasks{}).Where("is_complete = ? AND is_back_log = ?", false, false).Find(&tasks)

	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}

func (db *DB) GetTasksByCategory(categoryID uint) ([]Tasks, error) {
	var tasks []Tasks
	result := db.Where("category_id = ?", categoryID).Find(&tasks)
	return tasks, result.Error
}

// GetTasksDueSoon returns incomplete, non-backlog tasks whose due date falls
// within [now, now+window]. User is preloaded so callers can read Browser config.
func (db *DB) GetTasksDueSoon(window time.Duration) ([]Tasks, error) {
	var tasks []Tasks
	now := time.Now()
	result := db.Preload("User").
		Where("is_completed = ? AND is_back_log = ? AND due_date > ? AND due_date <= ?",
			false, false, now, now.Add(window)).
		Find(&tasks)
	return tasks, result.Error
}

func (db *DB) DeleteTask(TaskID uint) error {
	result := db.Delete(&Tasks{}, TaskID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
