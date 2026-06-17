package controllers

import (
	"ETM/pkg/app"
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

func IsAdminUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		App := c.MustGet("App").(*app.App)
		userID, ok := c.Get("userID")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}
		user, err := App.DB.GetUser(uint(userID.(float64)))
		if err != nil || user.IsAdmin != "true" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "admin access required"})
			return
		}
		c.Next()
	}
}
