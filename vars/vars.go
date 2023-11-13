// Global variales

package vars

import (
	"os"
)

var Token string
var Username string
var ConnectionString string
var Host string
var Port string

func Init() {
	Token = os.Getenv("TOKEN")
	Username = os.Getenv("USERNAME")
	ConnectionString = os.Getenv("CONNECTION_STRING")
	Host = os.Getenv("HOST")
	Port = os.Getenv("PORT")

}
