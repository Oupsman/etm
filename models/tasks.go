package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type Tasks struct {
	gorm.Model
	Name       string
	Comment    string
	IsComplete bool
	IsBackLog  bool
	CategoryID uint
	Category   Category
	Priority   bool
	Urgency    bool
	DueDate    time.Time
}

func GetTasks(c *gin.Context) {
	var db = ConnectToDb()
	var Tasks = []Tasks{}
	result := db.Find(&Tasks)
	if result.Error != nil {
		panic(result.Error)
	}
	c.JSON(http.StatusOK, Tasks)
}

func GetTask(c *gin.Context) {
	var db = ConnectToDb()
	var task = Tasks{}
	var id = c.Query("id")
	result := db.First(&task, id)
	if result.Error != nil {
		c.JSON(http.StatusForbidden, "Unable to get task ID ")
	}

	c.JSON(http.StatusOK, task)
}

func CreateTask(c *gin.Context) {

	var db = ConnectToDb()
	var name = c.Query("name")
	var comment = c.Query("comment")
	var dueDate time.Time
	var categoryid uint

	dueDate, _ = time.Parse(time.RFC3339, c.Query("duedate"))
	categID, _ := strconv.ParseUint(c.Query("categoryid"), 10, 32)
	categoryid = uint(categID)

	var task = Tasks{
		Name:       name,
		Comment:    comment,
		IsBackLog:  true,
		Priority:   false,
		Urgency:    false,
		CategoryID: categoryid,
		DueDate:    dueDate,
	}
	result := db.Create(&task)
	if result.Error != nil {
		c.JSON(http.StatusForbidden, task)
	}

	c.JSON(http.StatusOK, task)

}

func UpdateTask(c *gin.Context) {
	var db = ConnectToDb()
	var id int
	id, _ = strconv.Atoi(c.Query("id"))
	var name = c.Query("name")
	var comment = c.Query("comment")
	var dueDate time.Time
	var task = Tasks{}
	var priority bool
	var urgency bool
	var completed bool

	dueDate, _ = time.Parse(time.RFC3339, c.Query("duedate"))
	result := db.First(&task, id)
	if result.Error != nil {
		c.JSON(http.StatusForbidden, task)
	}

	task.Comment = comment
	task.Name = name
	task.DueDate = dueDate
	task.Priority = priority
	task.Urgency = urgency
	task.IsComplete = completed

	result = db.Save(&task)
	if result.Error != nil {
		c.JSON(http.StatusForbidden, task)
	}

	c.JSON(http.StatusOK, task)
}
