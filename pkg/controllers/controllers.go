package controllers

import (
	"ETM/pkg/app"
	"ETM/pkg/utils"
	"ETM/pkg/vars"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		if bearer := c.GetHeader("Authorization"); bearer != "" {
			parts := strings.SplitN(bearer, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid authorization format"})
				return
			}
			token := parts[1]

			if strings.HasPrefix(token, "etm_") {
				// API key auth
				if len(token) >= 12 {
					App := c.MustGet("App").(*app.App)
					device, err := App.DB.GetDeviceByPrefix(token[:12])
					if err == nil && utils.VerifyAPIKey(token, device.APIKeyHash) {
						go App.DB.UpdateDeviceLastUsed(device.ID)
						c.Set("userID", float64(device.UserID))
						c.Set("uuid", device.User.UUID.String())
						c.Next()
						return
					}
				}
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
				return
			}

			// JWT auth
			claims, err := utils.ParseToken(token)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
				return
			}
			c.Set("userID", claims["sub"])
			c.Set("uuid", claims["uuid"])
			c.Next()
			return
		}

		// No Authorization header — try session cookie
		cookie, err := c.Cookie("etm_session")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}
		claims, err := utils.ParseToken(cookie)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			return
		}
		// Slide the cookie if less than 5 minutes remain
		if exp, ok := claims["exp"].(float64); ok {
			if time.Until(time.Unix(int64(exp), 0)) < 5*time.Minute {
				if newToken := mintJWT(claims["sub"].(float64), claims["uuid"].(string)); newToken != "" {
					setSessionCookie(c, newToken)
				}
			}
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

// mintJWT creates a fresh JWT for the given user, valid for vars.TokenDuration.
func mintJWT(userID float64, userUUID string) string {
	claims := jwt.MapClaims{
		"authorized": true,
		"exp":        time.Now().Add(vars.TokenDuration).Unix(),
		"iss":        "etm",
		"sub":        userID,
		"uuid":       userUUID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(vars.SecretKey))
	if err != nil {
		return ""
	}
	return s
}

// setSessionCookie writes an HttpOnly session cookie to the response.
func setSessionCookie(c *gin.Context, token string) {
	secure := c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "etm_session",
		Value:    token,
		MaxAge:   1800,
		Path:     "/",
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

// clearSessionCookie removes the session cookie.
func clearSessionCookie(c *gin.Context) {
	secure := c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "etm_session",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

// getUUIDFromCtx extracts the user UUID set by IsAuthorized.
func getUUIDFromCtx(c *gin.Context) (uuid.UUID, bool) {
	val, ok := c.Get("uuid")
	if !ok {
		return uuid.UUID{}, false
	}
	u, err := uuid.Parse(val.(string))
	if err != nil {
		return uuid.UUID{}, false
	}
	return u, true
}

// getIDFromCtx extracts the numeric user ID set by IsAuthorized.
func getIDFromCtx(c *gin.Context) (uint, bool) {
	val, ok := c.Get("userID")
	if !ok {
		return 0, false
	}
	return uint(val.(float64)), true
}
