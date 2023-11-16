package models

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
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
	CategoryID := c.Query("categoryId")
	result := db.Find(&Tasks, CategoryID)
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
	// body, _ := io.ReadAll(c.Request.Body)
	// println(string(body))

	var db = ConnectToDb()

	type TaskBody struct {
		Name    string `json:"name"`
		Comment string `json:"comment"`
		DueDate string `json:"duedate"`
	}

	taskBody := TaskBody{}
	var dueDate time.Time

	err := c.ShouldBindJSON(&taskBody)
	switch {
	case errors.Is(err, io.EOF):
		fmt.Println("Error :", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing request body"})
		return
	case err != nil:
		fmt.Println("Error :", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dueDate, _ = time.Parse(time.RFC3339, taskBody.DueDate)
	categID, _ := strconv.ParseUint(c.Query("categoryid"), 10, 32)

	categoryid := uint(categID)

	var task = Tasks{
		Name:       taskBody.Name,
		Comment:    taskBody.Comment,
		IsBackLog:  true,
		Priority:   false,
		Urgency:    false,
		CategoryID: categoryid,
		DueDate:    dueDate,
	}

	fmt.Println("Name")
	fmt.Println(taskBody.Name)
	fmt.Println("Comment")
	fmt.Println(taskBody.Comment)
	result := db.Create(&task)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unable to add task to database"})
		return
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

func DeleteTask(c *gin.Context) {
	var db = ConnectToDb()

	var id int
	id, _ = strconv.Atoi(c.Query("id"))

	var Task = Tasks{}

	result := db.Find(&Task, id)
	if result.Error != nil {
		c.JSON(http.StatusForbidden, Task)
	}

	result = db.Delete(&Task)
	if result.Error != nil {
		c.JSON(http.StatusForbidden, Task)
	}
	c.JSON(http.StatusOK, Task)

}
