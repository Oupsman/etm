package controllers

import (
	"ETM/pkg/app"
	"ETM/pkg/utils"
	"ETM/pkg/vars"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

var (
	oidcProvider *oidc.Provider
	oidcVerifier *oidc.IDTokenVerifier
	oauth2Config oauth2.Config
)

// InitOIDC sets up the OIDC provider. Call once from main.go after App is created.
func InitOIDC() error {
	if vars.OIDCEnabled != "true" {
		return nil
	}
	ctx := context.Background()
	var err error
	oidcProvider, err = oidc.NewProvider(ctx, vars.OIDCIssuerURL)
	if err != nil {
		return err
	}
	oauth2Config = oauth2.Config{
		ClientID:     vars.OIDCClientID,
		ClientSecret: vars.OIDCClientSecret,
		RedirectURL:  vars.OIDCRedirectURL,
		Endpoint:     oidcProvider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	oidcVerifier = oidcProvider.Verifier(&oidc.Config{
		ClientID: vars.OIDCClientID,
	})
	return nil
}

// OIDCStatus returns whether OIDC login is available.
func OIDCStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"enabled": vars.OIDCEnabled == "true"})
}

// OIDCLogin redirects the browser to the provider's login page.
func OIDCLogin(c *gin.Context) {
	if vars.OIDCEnabled != "true" {
		c.JSON(http.StatusNotFound, gin.H{"error": "OIDC not enabled"})
		return
	}
	state := uuid.NewString()
	c.SetCookie("oidc_state", state, 300, "/", "", false, true)
	c.Redirect(http.StatusFound, oauth2Config.AuthCodeURL(state))
}

// OIDCLink starts the OIDC flow to attach an OIDC identity to an existing local
// account. The JWT is accepted as a query parameter because this is a browser
// redirect, not an XHR — the Authorization header cannot be set.
func OIDCLink(c *gin.Context) {
	if vars.OIDCEnabled != "true" {
		c.JSON(http.StatusNotFound, gin.H{"error": "OIDC not enabled"})
		return
	}

	token := c.Query("token")
	if token == "" {
		token = strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	}
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	userUUIDStr, _ := claims["uuid"].(string)
	userUUID, err := uuid.Parse(userUUIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token claims"})
		return
	}

	App := c.MustGet("App").(*app.App)
	user, err := App.DB.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if user.OIDCSubject != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "account is already linked to an OIDC identity"})
		return
	}

	state := uuid.NewString()
	c.SetCookie("oidc_state", state, 300, "/", "", false, true)
	c.SetCookie("oidc_link_uuid", user.UUID.String(), 300, "/", "", false, true)
	c.Redirect(http.StatusFound, oauth2Config.AuthCodeURL(state))
}

// OIDCUnlink removes the OIDC identity from the authenticated user's account.
func OIDCUnlink(c *gin.Context) {
	App := c.MustGet("App").(*app.App)
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := App.DB.UnlinkOIDC(uint(userID.(float64))); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OIDC account unlinked"})
}

// OIDCCallback handles the redirect back from the provider.
// It serves two purposes depending on whether an oidc_link_uuid cookie is set:
//   - Link flow: attaches the verified OIDC identity to the existing local user.
//   - Login flow: finds or creates a user, then issues a JWT.
func OIDCCallback(c *gin.Context) {
	if vars.OIDCEnabled != "true" {
		c.JSON(http.StatusNotFound, gin.H{"error": "OIDC not enabled"})
		return
	}

	savedState, err := c.Cookie("oidc_state")
	if err != nil || savedState != c.Query("state") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}
	c.SetCookie("oidc_state", "", -1, "/", "", false, true)

	ctx := context.Background()
	oauth2Token, err := oauth2Config.Exchange(ctx, c.Query("code"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token exchange failed"})
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no id_token in response"})
		return
	}
	idToken, err := oidcVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "id_token verification failed"})
		return
	}

	var claims struct {
		Sub               string `json:"sub"`
		Email             string `json:"email"`
		PreferredUsername string `json:"preferred_username"`
	}
	if err := idToken.Claims(&claims); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse claims"})
		return
	}

	App := c.MustGet("App").(*app.App)

	// Link flow: oidc_link_uuid cookie was set by OIDCLink.
	if linkUUID, cookieErr := c.Cookie("oidc_link_uuid"); cookieErr == nil && linkUUID != "" {
		c.SetCookie("oidc_link_uuid", "", -1, "/", "", false, true)

		parsedUUID, err := uuid.Parse(linkUUID)
		if err != nil {
			c.Redirect(http.StatusFound, "/auth/callback?error=invalid_link_state")
			return
		}
		targetUser, err := App.DB.GetUserByUUID(parsedUUID)
		if err != nil {
			c.Redirect(http.StatusFound, "/auth/callback?error=user_not_found")
			return
		}
		// Reject if this OIDC subject is already bound to a different account.
		if existing, err := App.DB.GetUserByOIDCSubject(claims.Sub); err == nil && existing.ID != targetUser.ID {
			c.Redirect(http.StatusFound, "/auth/callback?error=oidc_already_linked")
			return
		}
		if err := App.DB.LinkOIDCToUser(targetUser.ID, claims.Sub, "oidc"); err != nil {
			App.Logger.Error().Err(err).Msg("OIDCCallback: failed to link OIDC to user")
			c.Redirect(http.StatusFound, "/auth/callback?error=link_failed")
			return
		}
		c.Redirect(http.StatusFound, "/auth/callback?linked=true")
		return
	}

	// Normal login flow: find or create the local user.
	username := claims.PreferredUsername
	if username == "" {
		username = claims.Email
	}
	user, created, err := App.DB.FindOrCreateOIDCUser(claims.Sub, claims.Email, username)
	if err != nil {
		App.Logger.Error().Err(err).Str("sub", claims.Sub).Str("email", claims.Email).Msg("OIDC FindOrCreateOIDCUser error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not find or create user"})
		return
	}
	if created {
		if err := App.DB.SeedDemoContent(user); err != nil {
			App.Logger.Warn().Err(err).Uint("userID", user.ID).Msg("failed to seed demo content for new OIDC user")
		}
	}

	expirationTime := time.Now().Add(30 * time.Minute)
	jwtClaims := jwt.MapClaims{
		"authorized": true,
		"exp":        expirationTime.Unix(),
		"iss":        "etm",
		"sub":        user.ID,
		"uuid":       user.UUID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString([]byte(vars.SecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	c.Redirect(http.StatusFound, "/auth/callback?token="+tokenString)
}
