package controllers

import (
	"ETM/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Validate the JWT token
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

		// Parsing the token
		claims, err := utils.ParseToken(reqToken[1])
		if err != nil {
			switch err {
			case jwt.ErrSignatureInvalid:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token signature"})
			default:
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			}
			return
		}
		/*
			// 2. Getting and validating the deviceID
			deviceID := c.GetHeader("X-Device-ID")
			if deviceID == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Device ID is required"})
				return
			}

			// Verifying the deviceID
			tokenDeviceID, ok := claims["deviceID"].(string)
			if !ok || tokenDeviceID != deviceID {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid device"})
				return
			}
		*/
		// 3. Store info in context for further use
		c.Set("userID", claims["sub"])
		c.Set("uuid", claims["uuid"])

		c.Next()
	}
}
