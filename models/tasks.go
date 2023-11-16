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
	Name       string `json:"name"`
	Comment    string `json:"comment"`
	IsComplete bool   `json:"iscomplete"`
	IsBackLog  bool   `json:"isbacklog"`
	CategoryID uint   `json:"category-id"`
	Category   Category
	Priority   bool      `json:"priority"`
	Urgency    bool      `json:"urgency"`
	DueDate    time.Time `json:"duedate"`
}

type TaskBody struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Comment     string `json:"comment"`
	DueDate     string `json:"duedate"`
	IsBackLog   bool   `json:"isbacklog,omitempty"`
	IsCompleted bool   `json:"iscompleted,omitempty"`
	Priority    bool   `json:"priority,omitempty"`
	Urgency     bool   `json:"urgency,omitempty"`
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unable to get category ID"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func CreateTask(c *gin.Context) {
	// body, _ := io.ReadAll(c.Request.Body)
	// println(string(body))

	var db = ConnectToDb()

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

	dueDate, err = time.Parse(time.RFC3339, taskBody.DueDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	var dueDate time.Time

	taskBody := TaskBody{}
	var task = Tasks{}

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

	id, _ := strconv.Atoi(taskBody.Id)

	result := db.First(&task, id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error finding task in database"})
		return
	}

	task.Comment = taskBody.Comment
	task.Name = taskBody.Name
	task.DueDate = dueDate
	task.Priority = taskBody.Priority
	task.Urgency = taskBody.Urgency
	task.IsComplete = taskBody.IsCompleted
	task.IsBackLog = taskBody.IsBackLog

	fmt.Println(task.Name)
	fmt.Println(task.Comment)
	fmt.Println(task.Priority)
	fmt.Println(task.Urgency)

	result = db.Save(&task)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error updating task in database"})
		return
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
