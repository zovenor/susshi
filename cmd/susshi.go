package main

import (
	"log"

	"github.com/zovenor/susshi/entities/app"
	sshserver "github.com/zovenor/susshi/entities/ssh-server"
)

func main() {
	servers := sshserver.Servers{}

	a := app.New(&servers)
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
