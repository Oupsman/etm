package utils

import (
	"ETM/pkg/vars"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CompareHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ParseToken(tokenString string) (claims jwt.MapClaims, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(vars.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims = token.Claims.(jwt.MapClaims)
	return claims, nil
}

// ParseTokenAllowExpired verifies the token signature but does not reject it
// when only the exp claim has passed. Used exclusively by the refresh endpoint
// so that an expired-but-authentic token can still yield a new one.
func ParseTokenAllowExpired(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(vars.SecretKey), nil
	})

	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		// Allow the token only when the sole problem is expiry.
		// Any other error (bad signature, malformed, alg:none…) is still rejected.
		if !ok || ve.Errors&^jwt.ValidationErrorExpired != 0 {
			return nil, err
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims, nil
}

func GetUserID(tokenString string) (uint, error) {
	reqToken := strings.Split(tokenString, " ")[1]

	claims, err := ParseToken(reqToken)
	if err != nil {
		return 0, err
	}
	return uint(claims["sub"].(float64)), nil
}

func GetUserUUID(tokenString string) (uuid.UUID, error) {
	reqToken := strings.Split(tokenString, " ")[1]

	claims, err := ParseToken(reqToken)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.Parse(claims["uuid"].(string))
}

// GenerateAPIKey creates a new device API key.
// Returns the plaintext key (shown once to the user), a short prefix (stored
// in plain text for fast DB lookup), and the SHA-256 hash (stored in the DB).
func GenerateAPIKey() (plaintext, prefix, hash string, err error) {
	buf := make([]byte, 32)
	if _, err = rand.Read(buf); err != nil {
		return
	}
	plaintext = "etm_" + base64.RawURLEncoding.EncodeToString(buf)
	prefix = plaintext[:12]
	sum := sha256.Sum256([]byte(plaintext))
	hash = hex.EncodeToString(sum[:])
	return
}

// VerifyAPIKey reports whether plaintext hashes to storedHash.
func VerifyAPIKey(plaintext, storedHash string) bool {
	sum := sha256.Sum256([]byte(plaintext))
	expected := hex.EncodeToString(sum[:])
	return subtle.ConstantTimeCompare([]byte(expected), []byte(storedHash)) == 1
}
