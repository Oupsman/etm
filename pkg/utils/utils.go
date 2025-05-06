package utils

import (
	"ETM/pkg/vars"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"strings"
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
		return []byte(vars.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims = token.Claims.(jwt.MapClaims)
	/*if !ok {
		return nil, err
	}*/

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
