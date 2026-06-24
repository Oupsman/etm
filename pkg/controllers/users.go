package controllers

import (
	"ETM/pkg/app"
	"ETM/pkg/crypto"
	"ETM/pkg/models"
	"ETM/pkg/types"
	"ETM/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
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

	tokenString := mintJWT(float64(existingUser.ID), existingUser.UUID.String())
	if tokenString == "" {
		c.JSON(500, gin.H{"error": "could not generate token"})
		return
	}

	setSessionCookie(c, tokenString)
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
	clearSessionCookie(c)
	c.JSON(200, gin.H{"success": "user logged out"})
}

func RefreshToken(c *gin.Context) {
	userID, ok := getIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uuidVal, _ := c.Get("uuid")

	newToken := mintJWT(float64(userID), uuidVal.(string))
	if newToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	// Refresh the session cookie if the request came via cookie auth
	if _, err := c.Cookie("etm_session"); err == nil {
		setSessionCookie(c, newToken)
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken})
}

func WhoAmI(c *gin.Context) {
	userID, ok := getIDFromCtx(c)
	if !ok {
		c.JSON(400, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(200, gin.H{"sub": userID})
}

func GetUser(c *gin.Context) {
	App := c.MustGet("App")
	db := App.(*app.App).DB

	userUUID, ok := getUUIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := db.GetUserByUUID(userUUID)
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

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userUUID, ok := getUUIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	updatedUser.UUID = userUUID
	currentUser, err := db.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if !utils.CompareHashPassword(user.OldPassword, currentUser.Password) {
		c.JSON(401, gin.H{"error": "wrong password"})
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

	if err := db.UpdateUser(updatedUser); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "user updated successfully"})
}

func UpdateUserSubscription(c *gin.Context) {
	var requestBody types.BrowserConfig

	App := c.MustGet("App")
	db := App.(*app.App).DB

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userUUID, ok := getUUIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	currentUser, err := db.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(400, gin.H{"error": "user not found"})
		return
	}

	currentUser.Browser = crypto.EncryptedString(requestBody.Subscription)
	if err := db.UpdateUser(currentUser); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "user updated successfully"})
}

func SendTestNotification(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	userUUID, ok := getUUIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
	if err := BrowserSend(msg, string(user.Browser)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "test notification sent"})
}

func UpdateUserPreferences(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	userUUID, ok := getUUIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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

	userUUID, ok := getUUIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := db.GetUserByUUID(userUUID)
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

	userID, ok := getIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	device.UserID = userID
	if err := db.CreateDevice(device); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create device"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": "device created"})
}
