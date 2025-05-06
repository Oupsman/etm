package controllers

import (
	"ETM/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

func IsAuthorized() gin.HandlerFunc {

	return func(c *gin.Context) {

		bearerToken := c.GetHeader("Authorization")
		reqToken := strings.Split(bearerToken, " ")[1]
		_, err := utils.ParseToken(reqToken)

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized",
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "bad request",
			})
			return
		}

		c.Next()
	}
}
