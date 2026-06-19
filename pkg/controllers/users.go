package controllers

import (
	"ETM/pkg/app"
	"ETM/pkg/models"
	"ETM/pkg/types"
	"ETM/pkg/utils"
	"ETM/pkg/vars"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Login(c *gin.Context) {

	App := c.MustGet("App")
	db := App.(*app.App).DB

	var creds types.UserBody

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.Users

	db.Where("username = ?", creds.Username).First(&existingUser)

	if existingUser.ID == 0 {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}

	if !utils.CompareHashPassword(creds.Password, existingUser.Password) {
		c.JSON(400, gin.H{"error": "invalid password"})
		return
	}

	expirationTime := time.Now().Add(30 * time.Minute)
	/*
		// Getting the deviceID
		deviceID := c.GetHeader("X-Device-ID")
		if deviceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Device ID is required"})
			return
		}
	*/
	// Create JWT
	claims := jwt.MapClaims{
		"authorized": true,
		"exp":        expirationTime.Unix(),
		"iss":        "etm",
		"sub":        existingUser.ID,
		"uuid":       existingUser.UUID.String(),
		//		"deviceid":   deviceID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(vars.SecretKey))

	if err != nil {
		c.JSON(500, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(200, gin.H{"token": tokenString})
}

func Register(c *gin.Context) {
	App := c.MustGet("App")
	db := App.(*app.App).DB

	var user types.UserBody

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.Users

	db.Where("username = ?", user.Username).First(&existingUser)

	if existingUser.ID != 0 {
		c.JSON(409, gin.H{"error": "user already exists"})
		return
	}

	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(user.Password)

	if errHash != nil {
		c.JSON(500, gin.H{"error": "could not generate password hash"})
		return
	}

	count, err := db.CountUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": "could not count users"})
		return
	}
	isFirst := count == 0

	newUser, err := db.CreateUser(user, isFirst)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not create user"})
		return
	}

	if isFirst {
		if _, err := db.CreateGroup("Administrators", newUser.ID); err != nil {
			App.(*app.App).Logger.Warn().Err(err).Msg("failed to create admin group for first user")
		}
	}

	if err := db.SeedDemoContent(newUser); err != nil {
		App.(*app.App).Logger.Warn().Err(err).Uint("userID", newUser.ID).Msg("failed to seed demo content for new user")
	}

	c.JSON(201, gin.H{"success": "user registered"})
}

func ListAllUsers(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	users, err := App.DB.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func SearchUsersHandler(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	query := c.Query("q")
	if len(query) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query must be at least 2 characters"})
		return
	}
	users, err := App.DB.SearchUsers(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func Logout(c *gin.Context) {

	c.JSON(200, gin.H{"success": "user logged out"})
}

func RefreshToken(c *gin.Context) {
	bearerToken := c.Request.Header.Get("Authorization")
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid authorization header"})
		return
	}

	claims, err := utils.ParseToken(parts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	expirationTime := time.Now().Add(30 * time.Minute)

	newClaims := jwt.MapClaims{
		"authorized": true,
		"iss":        "etm",
		"sub":        claims["sub"],
		"uuid":       claims["uuid"],
		"exp":        expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenString, err := token.SignedString([]byte(vars.SecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func WhoAmI(c *gin.Context) {
	bearerToken := c.Request.Header.Get("Authorization")
	UserID, err := utils.GetUserID(bearerToken)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"sub": UserID})
}

func GetUser(c *gin.Context) {
	App := c.MustGet("App")
	db := App.(*app.App).DB

	bearerToken := c.Request.Header.Get("Authorization")
	UserUUID, err := utils.GetUserUUID(bearerToken)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := db.GetUserByUUID(UserUUID)

	if err != nil {
		c.JSON(400, gin.H{"error": "user not found"})
		return
	}
	c.JSON(200, user)
}

func UpdateUser(c *gin.Context) {

	var updatedUser models.Users
	var user types.UserBody

	App := c.MustGet("App")
	db := App.(*app.App).DB

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	bearerToken := c.Request.Header.Get("Authorization")
	UserUUID, err := utils.GetUserUUID(bearerToken)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	updatedUser.UUID = UserUUID
	currentUser, err := db.GetUserByUUID(UserUUID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if !utils.CompareHashPassword(user.OldPassword, currentUser.Password) {
		c.JSON(401, gin.H{
			"error": "wrong password",
		})
		return
	}

	updatedUser.Email = user.Email
	updatedUser.Username = currentUser.Username
	if user.Password != "" {
		newHash, err := utils.GenerateHashPassword(user.Password)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		updatedUser.Password = newHash
	} else {
		updatedUser.Password = currentUser.Password
	}

	err = db.UpdateUser(updatedUser)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user updated successfully"})
}

func UpdateUserSubscription(c *gin.Context) {

	var requestBody types.BrowserConfig

	App := c.MustGet("App")
	db := App.(*app.App).DB

	err := c.ShouldBindJSON(&requestBody)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	bearerToken := c.Request.Header.Get("Authorization")
	UserUUID, err := utils.GetUserUUID(bearerToken)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser, err := db.GetUserByUUID(UserUUID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser.Browser = requestBody.Subscription

	err = db.UpdateUser(currentUser)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "user updated successfully"})

}

func SendTestNotification(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	bearerToken := c.Request.Header.Get("Authorization")
	userUUID, err := utils.GetUserUUID(bearerToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	user, err := db.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}
	if user.Browser == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no push subscription found — enable notifications in your profile first"})
		return
	}
	msg := `{"title":"ETM test","body":"Push notifications are working!"}`
	if err := BrowserSend(msg, user.Browser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "test notification sent"})
}

func UpdateUserPreferences(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	bearerToken := c.Request.Header.Get("Authorization")
	userUUID, err := utils.GetUserUUID(bearerToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var prefs types.UserPreferencesBody
	if err := c.ShouldBindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := db.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.ActiveCategoryID = prefs.ActiveCategoryID
	if err := db.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "preferences updated"})
}

func GetUserDevices(c *gin.Context) {
	App := c.MustGet("App")
	db := App.(*app.App).DB
	bearerToken := c.Request.Header.Get("Authorization")
	UserUUID, err := utils.GetUserUUID(bearerToken)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := db.GetUserByUUID(UserUUID)

	if err != nil {
		c.JSON(400, gin.H{"error": "user not found"})
		return
	}
	c.JSON(200, user)
}

func CreateUserDevice(c *gin.Context) {
	App := c.MustGet("App")
	db := App.(*app.App).DB

	var device types.DeviceBody
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bearerToken := c.Request.Header.Get("Authorization")
	userID, err := utils.GetUserID(bearerToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	device.UserID = userID

	err = db.CreateDevice(device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create device"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": "device created"})
}
