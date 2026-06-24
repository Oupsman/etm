package controllers

import (
	"ETM/pkg/app"
	"ETM/pkg/crypto"
	models2 "ETM/pkg/models"
	"ETM/pkg/types"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	categoryID, err := strconv.Atoi(c.Param("categoryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	userUUID, ok := getUUIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user, err := db.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !db.CanUserViewCategory(user.ID, uint(categoryID)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	tasks, err := db.GetTasksByCategory(uint(categoryID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unable to list tasks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func GetTask(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	taskID, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	userUUID, ok := getUUIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user, err := db.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := db.GetTask(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if !db.CanUserViewCategory(user.ID, task.CategoryID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func CreateTask(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	userUUID, ok := getUUIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user, err := db.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskBody := types.TaskBody{}
	err = c.ShouldBindJSON(&taskBody)
	switch {
	case errors.Is(err, io.EOF):
		App.Logger.Error().Err(err).Msg("CreateTask: missing request body")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing request body"})
		return
	case err != nil:
		App.Logger.Error().Err(err).Msg("CreateTask: failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !db.CanUserEditCategory(user.ID, taskBody.CategoryID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	dueDate, err := time.Parse(time.RFC3339, taskBody.DueDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid due date format, expected RFC3339"})
		return
	}

	task := models2.Tasks{
		Name:       crypto.EncryptedString(taskBody.Name),
		Comment:    crypto.EncryptedString(taskBody.Comment),
		Link:       crypto.EncryptedString(taskBody.Link),
		IsBackLog:  true,
		Priority:   false,
		Urgency:    false,
		CategoryID: taskBody.CategoryID,
		DueDate:    dueDate,
		UserID:     user.ID,
	}

	if err := db.CreateTask(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"task": task})
}

func UpdateTask(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	userUUID, ok := getUUIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user, err := db.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskID, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	taskBody := types.TaskBody{}
	err = c.ShouldBindJSON(&taskBody)
	switch {
	case errors.Is(err, io.EOF):
		App.Logger.Error().Err(err).Msg("UpdateTask: missing request body")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing request body"})
		return
	case err != nil:
		App.Logger.Error().Err(err).Msg("UpdateTask: failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := db.GetTask(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if !db.CanUserEditCategory(user.ID, task.CategoryID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	dueDate, err := time.Parse(time.RFC3339, taskBody.DueDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid due date format, expected RFC3339"})
		return
	}
	if taskBody.CategoryID != 0 && taskBody.CategoryID != task.CategoryID {
		if !db.CanUserEditCategory(user.ID, taskBody.CategoryID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied to target category"})
			return
		}
		task.CategoryID = taskBody.CategoryID
	}

	task.Comment = crypto.EncryptedString(taskBody.Comment)
	task.Name = crypto.EncryptedString(taskBody.Name)
	task.Link = crypto.EncryptedString(taskBody.Link)
	task.DueDate = dueDate
	task.Priority = taskBody.Priority
	task.Urgency = taskBody.Urgency
	task.IsCompleted = taskBody.IsCompleted
	task.IsStarted = taskBody.IsStarted
	task.IsBackLog = taskBody.IsBackLog

	if err := db.UpdateTask(task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": task})
}

func DeleteTask(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	userUUID, ok := getUUIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user, err := db.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if !db.CanUserEditCategory(user.ID, task.CategoryID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := db.DeleteTask(task.ID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not delete task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}

func CheckTasks() error {
	return nil
}
