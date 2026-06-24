package controllers

import (
	"ETM/pkg/app"
	"ETM/pkg/types"
	"ETM/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RegisterAPIDevice creates a new device and returns the API key once.
// POST /api/v1/devices
func RegisterAPIDevice(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	userID, ok := getIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body types.DeviceRegisterBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plaintext, prefix, hash, err := utils.GenerateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate API key"})
		return
	}

	device, err := db.RegisterDevice(userID, body.DeviceName, prefix, hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not register device"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"device":  device,
		"api_key": plaintext,
	})
}

// ListAPIDevices returns all devices for the authenticated user.
// GET /api/v1/devices
func ListAPIDevices(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	userID, ok := getIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	devices, err := db.GetDevices(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve devices"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"devices": devices})
}

// DeleteAPIDevice revokes a device by ID.
// DELETE /api/v1/devices/:id
func DeleteAPIDevice(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	userID, ok := getIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	device, err := db.GetDeviceByID(uint(deviceID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}
	if device.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := db.DeleteDevice(uint(deviceID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete device"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "device deleted"})
}
