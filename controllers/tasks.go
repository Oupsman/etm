package controllers

import (
	"ETM/models"
	"ETM/types"
	"ETM/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"time"
)

func GetTasks(c *gin.Context) {

	CategoryID, err := strconv.Atoi(c.Param("CategoryId"))
	bearerToken := c.Request.Header.Get("Authorization")
	UserID, err := utils.GetUserID(bearerToken)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tasks, err := models.GetTasks(UserID, CategoryID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unable to list tasks"})
		return
	}

	c.JSON(200, tasks)
}

func GetTask(c *gin.Context) {
	TaskID, err := strconv.Atoi(c.Param("taskId"))
	bearerToken := c.Request.Header.Get("Authorization")
	UserID, err := utils.GetUserID(bearerToken)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	task, err := models.GetTask(TaskID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if task.UserID != UserID {
		c.JSON(403, gin.H{"error": "You do not have access to this task"})
	}
	c.JSON(200, task)
}

func CreateTask(c *gin.Context) {

	taskBody := types.TaskBody{}
	var dueDate time.Time
	bearerToken := c.Request.Header.Get("Authorization")
	UserID, err := utils.GetUserID(bearerToken)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err = c.ShouldBindJSON(&taskBody)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

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

	CategoryID, _ := strconv.Atoi(taskBody.CategoryID)

	var task = models.Tasks{
		Name:       taskBody.Name,
		Comment:    taskBody.Comment,
		IsBackLog:  true,
		Priority:   false,
		Urgency:    false,
		CategoryID: uint(CategoryID),
		DueDate:    dueDate,
		UserID:     UserID,
	}

	err = models.CreateTask(task)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{
		"task": task,
	})
}

func UpdateTask(c *gin.Context) {
	var dueDate time.Time

	taskBody := types.TaskBody{}

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

	task, err := models.GetTask(taskBody.Id)
	if err != nil {
		c.JSON(400, gin.H{"error": "task not found"})
		return
	}

	task.Comment = taskBody.Comment
	task.Name = taskBody.Name
	task.DueDate = dueDate
	task.Priority = taskBody.Priority
	task.Urgency = taskBody.Urgency
	task.IsComplete = taskBody.IsCompleted
	task.IsBackLog = taskBody.IsBackLog

	err = models.UpdateTask(task)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"task": task,
	})
}

func DeleteTask(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("taskId"))

	Task, err := models.GetTask(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	result := models.Db.Delete(&Task)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "error deleting task from database" + result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "task deleted successfully",
	})
}

func CheckTasks() error {

	return nil
}
