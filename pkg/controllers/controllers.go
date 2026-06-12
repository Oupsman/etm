package controllers

import (
	"ETM/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		if bearerToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is required"})
			return
		}

		reqToken := strings.Split(bearerToken, " ")
		if len(reqToken) != 2 || reqToken[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid authorization format"})
			return
		}

		claims, err := utils.ParseToken(reqToken[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			return
		}

		c.Set("userID", claims["sub"])
		c.Set("uuid", claims["uuid"])

		c.Next()
	}
}
