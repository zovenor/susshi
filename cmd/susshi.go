package main

import (
	"log"
	"os"
	"path"

	"github.com/zovenor/susshi/internal/app"
	sshserver "github.com/zovenor/susshi/internal/ssh-server"
)

const serversPath = "/.susshi/servers.json"

func main() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	servers := sshserver.Servers{}
	err = servers.ImportFromFile(path.Join(homePath, serversPath))
	if err != nil {
		log.Fatal(err)
	}
	a := app.New(&servers)
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
