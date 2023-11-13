package main

import (
	"ETM/vars"
	"fmt"
	_ "gorm.io/driver/sqlite"
	"log"
	"net"
)

func main() {

	vars.Init()

	fmt.Println("Starting Eisenhower Task Manager")
	fmt.Printf("Username: %s\n", vars.Username)
	fmt.Printf("Token: %s\n", vars.Token)
	fmt.Printf("Connection String: %s\n", vars.ConnectionString)
	fmt.Printf("Listening on %s:%s\n", vars.Host, vars.Port)

	addr := net.JoinHostPort(vars.Host, vars.Port)
	if err := runHttp(addr); err != nil {
		log.Fatal(err)
	}

}
