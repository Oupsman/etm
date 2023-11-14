// Global variales

package vars

import (
	"os"
)

var Host string
var Port string

func Init() {
	Host = os.Getenv("HOST")
	Port = os.Getenv("PORT")
}
