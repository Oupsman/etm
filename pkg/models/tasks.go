package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Tasks struct {
	gorm.Model
	Name       string `json:"name"`
	Comment    string `json:"comment"`
	IsComplete bool   `json:"iscomplete"`
	IsBackLog  bool   `json:"isbacklog"`
	CategoryID uint   `json:"category-id"`
	Category   Category
	Priority   bool      `json:"priority"`
	Urgency    bool      `json:"urgency"`
	DueDate    time.Time `json:"duedate"`
	User       Users
	UserID     uint `json:"userid"`
}

func (db *DB) GetTasks(UserUUID uuid.UUID, CategoryId int) ([]Tasks, error) {
	// var Db = ConnectToDb()
	var Tasks = []Tasks{}
	result := db.Where("user_uuid = ? AND category_id = ?", UserUUID, CategoryId).Find(&Tasks)
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

func (db *DB) CreateTask(task Tasks) error {

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
