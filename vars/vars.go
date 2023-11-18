// Global variales

package vars

import (
	"os"
)

var Host string
var Port string

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
}
