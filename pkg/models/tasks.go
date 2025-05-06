package models

import (
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

func GetTasks(UserID uint, CategoryId int) ([]Tasks, error) {
	// var Db = ConnectToDb()
	var Tasks = []Tasks{}
	result := Db.Where("user_id = ? AND category_id = ?", UserID, CategoryId).Find(&Tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return Tasks, nil
}

func GetTask(TaskID int) (*Tasks, error) {
	// var Db = ConnectToDb()
	var task = Tasks{}
	result := Db.First(&task, TaskID)

	if result.Error != nil {

		return nil, result.Error
	}

	return &task, nil
}

func CreateTask(task Tasks) error {

	result := Db.Create(&task)
	if result.Error != nil {

		return result.Error
	}

	return nil
}

func UpdateTask(task *Tasks) error {
	// var Db = ConnectToDb()

	result := Db.Save(&task)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetActiveTasks() ([]Tasks, error) {

	var tasks []Tasks

	result := Db.Model(&Tasks{}).Where("is_complete = ? AND is_back_log = ?", false, false).Find(&tasks)

	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}
