package controllers

import (
	"ETM/models"
	"ETM/types"
	"ETM/utils"
	"ETM/vars"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"strings"
	"time"
)

func Login(c *gin.Context) {

	var user models.Users

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.Users

	models.Db.Where("name = ?", user.Name).First(&existingUser)

	if existingUser.ID == 0 {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}

	errHash := utils.CompareHashPassword(user.Password, existingUser.Password)

	if !errHash {
		c.JSON(400, gin.H{"error": "invalid password"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"authorized": true,
		"exp":        expirationTime.Unix(),
		"iss":        "switchdb",
		"sub":        existingUser.ID,
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

	var user models.Users

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.Users

	models.Db.Where("Username = ?", user.Name).First(&existingUser)

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

	models.Db.Create(&user)

	c.JSON(201, gin.H{"success": "user registered"})
}

func Logout(c *gin.Context) {

	c.JSON(200, gin.H{"success": "user logged out"})
}

func RefreshToken(c *gin.Context) {
	bearerToken := c.Request.Header.Get("Authorization")
	reqToken := strings.Split(bearerToken, " ")[1]

	claims, err := utils.ParseToken(reqToken)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	newClaims := jwt.MapClaims{

		"authorized": true,
		"role":       claims["role"],
		"iss":        "switchdb",
		"sub":        claims["sub"],
		"exp":        expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenString, err := token.SignedString([]byte(vars.SecretKey))

	if err != nil {
		c.JSON(500, gin.H{"error": "could not generate token"})
		return
	}

	c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)
	c.JSON(200, gin.H{"token": tokenString})
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

	bearerToken := c.Request.Header.Get("Authorization")
	UserID, err := utils.GetUserID(bearerToken)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUser(UserID)

	if err != nil {
		c.JSON(400, gin.H{"error": "user not found"})
	}
	c.JSON(200, user)
}

func UpdateUser(c *gin.Context) {

	var updatedUser models.Users
	var user types.UserBody
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	bearerToken := c.Request.Header.Get("Authorization")
	UserID, err := utils.GetUserID(bearerToken)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	updatedUser.ID = uint(UserID)
	currentUser, err := models.GetUser(UserID)
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

	updatedUser.ID = uint(UserID)
	updatedUser.Email = user.Email
	updatedUser.Name = currentUser.Name
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

	err = models.UpdateUser(updatedUser)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user updated successfully"})
}
