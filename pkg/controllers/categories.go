package controllers

import (
	"ETM/pkg/app"
	"ETM/pkg/types"
	"ETM/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	bearerToken := c.Request.Header.Get("Authorization")
	userUUID, err := utils.GetUserUUID(bearerToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := db.GetUserByUUID(userUUID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	categories, err := db.GetCategories(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

func CreateCategory(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	bearerToken := c.Request.Header.Get("Authorization")
	userUUID, err := utils.GetUserUUID(bearerToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := db.GetUserByUUID(userUUID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var categoryBody types.CategoryBody
	if err := c.BindJSON(&categoryBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := db.CreateCategory(categoryBody, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if categoryBody.GroupID > 0 {
		if !db.IsGroupMember(categoryBody.GroupID, user.ID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you are not a member of that group"})
			return
		}
		if err := db.ShareCategoryWithGroup(category.ID, categoryBody.GroupID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "category created but sharing failed: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"category": category})
}

func ShareCategory(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	categoryID, err := strconv.Atoi(c.Param("categoryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	var body types.CategoryShareBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bearerToken := c.Request.Header.Get("Authorization")
	userUUID, err := utils.GetUserUUID(bearerToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	user, err := db.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat, err := db.GetCategory(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	if cat.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the category owner can share it"})
		return
	}

	if err := db.ShareCategoryWithGroup(uint(categoryID), body.GroupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "category shared"})
}

func UnshareCategory(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	categoryID, err := strconv.Atoi(c.Param("categoryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}
	groupID, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	bearerToken := c.Request.Header.Get("Authorization")
	userUUID, err := utils.GetUserUUID(bearerToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	user, err := db.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat, err := db.GetCategory(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	if cat.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the category owner can unshare it"})
		return
	}

	if err := db.UnshareCategoryFromGroup(uint(categoryID), uint(groupID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "category unshared"})
}

func GetCategoryShares(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	categoryID, err := strconv.Atoi(c.Param("categoryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	bearerToken := c.Request.Header.Get("Authorization")
	userUUID, err := utils.GetUserUUID(bearerToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	user, err := db.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat, err := db.GetCategory(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	if cat.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the category owner can view shares"})
		return
	}

	groups, err := db.GetCategoryShares(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"groups": groups})
}
