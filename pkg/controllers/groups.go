package controllers

import (
	"ETM/pkg/app"
	"ETM/pkg/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateGroup(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	var body types.GroupBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	group, err := db.CreateGroup(body.Name, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"group": group})
}

func GetGroups(c *gin.Context) {
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

	groups, err := db.GetGroupsByUser(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

func GetGroup(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	groupID, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
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

	if !db.IsGroupMember(uint(groupID), user.ID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member of this group"})
		return
	}

	group, err := db.GetGroup(uint(groupID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	members, err := db.GetGroupMembers(uint(groupID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"group": group, "members": members})
}

func DeleteGroup(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	groupID, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
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

	if !db.IsGroupOwner(uint(groupID), user.ID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the group owner can delete it"})
		return
	}

	if err := db.DeleteGroup(uint(groupID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "group deleted"})
}

func AddGroupMember(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	groupID, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	var body types.GroupMemberBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	if !db.IsGroupOwner(uint(groupID), user.ID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the group owner can add members"})
		return
	}

	if err := db.AddGroupMember(uint(groupID), body.UserID, body.Role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "member added"})
}

func UpdateGroupMember(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	groupID, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}
	memberID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var body types.GroupMemberBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	if !db.IsGroupOwner(uint(groupID), user.ID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the group owner can change roles"})
		return
	}

	if err := db.UpdateGroupMember(uint(groupID), uint(memberID), body.Role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "member updated"})
}

func RemoveGroupMember(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	db := App.DB

	groupID, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}
	memberID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
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

	if !db.IsGroupOwner(uint(groupID), user.ID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the group owner can remove members"})
		return
	}

	group, err := db.GetGroup(uint(groupID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}
	if group.OwnerID == uint(memberID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot remove the group owner"})
		return
	}

	if err := db.RemoveGroupMember(uint(groupID), uint(memberID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "member removed"})
}
