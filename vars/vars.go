// Global variales

package vars

import (
	"os"
)

var Host string
var Port string
var DbHost string
var DbPort string
var Database string
var Username string
var Password string

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func Init() {
	Host = getEnv("HOST", "127.0.0.1")
	Port = getEnv("PORT", "8080")
	DbHost = getEnv("DB_HOST", "127.0.0.1")
	DbPort = getEnv("DB_PORT", "5432")
	Database = getEnv("DATABASE", "etm")
	Username = getEnv("DB_USERNAME", "etmusr")
	Password = getEnv("DB_PASSWORD", "etmpass")

}
