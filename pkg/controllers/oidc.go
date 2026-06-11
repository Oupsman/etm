package controllers

import (
	"ETM/pkg/app"
	"ETM/pkg/vars"
	"context"
	"net/http"
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

// OIDCLogin redirects the browser to the provider's login page.
func OIDCLogin(c *gin.Context) {
	if vars.OIDCEnabled != "true" {
		c.JSON(http.StatusNotFound, gin.H{"error": "OIDC not enabled"})
		return
	}
	// Use a random UUID as the state nonce, stored in a short-lived cookie
	state := uuid.NewString()
	c.SetCookie("oidc_state", state, 300, "/", "", false, true)
	c.Redirect(http.StatusFound, oauth2Config.AuthCodeURL(state))
}

// OIDCCallback handles the redirect back from the provider.
func OIDCCallback(c *gin.Context) {
	if vars.OIDCEnabled != "true" {
		c.JSON(http.StatusNotFound, gin.H{"error": "OIDC not enabled"})
		return
	}

	// Validate state to prevent CSRF
	savedState, err := c.Cookie("oidc_state")
	if err != nil || savedState != c.Query("state") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}
	c.SetCookie("oidc_state", "", -1, "/", "", false, true)

	// Exchange the authorization code for tokens
	ctx := context.Background()
	oauth2Token, err := oauth2Config.Exchange(ctx, c.Query("code"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token exchange failed"})
		return
	}

	// Extract and verify the ID token
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

	// Parse standard claims
	var claims struct {
		Sub               string `json:"sub"`
		Email             string `json:"email"`
		PreferredUsername string `json:"preferred_username"`
	}
	if err := idToken.Claims(&claims); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse claims"})
		return
	}

	// Use preferred_username, fall back to email
	username := claims.PreferredUsername
	if username == "" {
		username = claims.Email
	}

	// Find or create the local user
	App := c.MustGet("App").(*app.App)
	user, err := App.DB.FindOrCreateOIDCUser(claims.Sub, claims.Email, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not find or create user"})
		return
	}

	// Issue a JWT identical in shape to the one Login() issues
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

	// Redirect the frontend with the token in the URL fragment
	// Adjust the path to match your Vue router's callback route
	c.Redirect(http.StatusFound, "/#/auth/callback?token="+tokenString)
}
