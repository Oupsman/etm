package vars

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Host string
var Port string
var DbHost string
var DbPort string
var Database string
var Username string
var Password string
var SecretKey string
var Dsn string
var OIDCEnabled string
var OIDCIssuerURL string
var OIDCClientID string
var OIDCClientSecret string
var OIDCRedirectURL string
var AllowedOrigin string
var VAPIDSubscriber string
var EncryptionKey string

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func Init() {
	if _, exists := os.LookupEnv("DB_HOST"); !exists {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, relying on system environment variables")
		}
	}

	Host = getEnv("HOST", "")
	Port = getEnv("PORT", "8080")
	DbHost = getEnv("DB_HOST", "127.0.0.1")
	DbPort = getEnv("DB_PORT", "5432")
	Database = getEnv("DATABASE", "etm")
	Username = getEnv("DB_USERNAME", "etm")
	Password = getEnv("DB_PASSWORD", "etmpass")
	SecretKey = getEnv("SECRET_KEY", "")
	OIDCEnabled = getEnv("OIDC_ENABLED", "false")
	OIDCIssuerURL = getEnv("OIDC_ISSUER_URL", "")
	OIDCClientID = getEnv("OIDC_CLIENT_ID", "")
	OIDCClientSecret = getEnv("OIDC_CLIENT_SECRET", "")
	OIDCRedirectURL = getEnv("OIDC_REDIRECT_URL", "")
	AllowedOrigin = getEnv("ALLOWED_ORIGIN", "*")
	VAPIDSubscriber = getEnv("VAPID_SUBSCRIBER", "")
	EncryptionKey = getEnv("ENCRYPTION_KEY", "")
	Dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", DbHost, Username, Password, Database, DbPort)

	if SecretKey == "" {
		log.Println("WARNING: SECRET_KEY is not set — JWT tokens will use an empty secret and can be trivially forged")
	}
}
