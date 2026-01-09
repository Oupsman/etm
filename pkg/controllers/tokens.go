package controllers

import (
	"ETM/pkg/app"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserTokens(c *gin.Context) {
	App := c.MustGet("App")
	db := App.(*app.App).DB
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non identifié"})
		return
	}

	uid := uint(userID.(float64))

	tokens, err := db.GetTokensByUserID(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des jetons"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func GetTokenByUUID(c *gin.Context) {
	App := c.MustGet("App")
	db := App.(*app.App).DB
	uuidParam := c.Param("uuid")
	tokenUUID, err := uuid.Parse(uuidParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format d'UUID invalide"})
		return
	}

	token, err := db.GetTokenByUUID(tokenUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jeton introuvable"})
		return
	}

	c.JSON(http.StatusOK, token)
}

func DeleteToken(c *gin.Context) {
	App := c.MustGet("App")
	db := App.(*app.App).DB
	idParam := c.Param("id")
	tokenID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	err = db.DeleteToken(uint(tokenID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du jeton"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jeton supprimé avec succès"})
}
