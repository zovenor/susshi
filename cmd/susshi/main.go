package main

import (
	"log"
	"os/exec"

	"github.com/zovenor/susshi/internal/app"
	"github.com/zovenor/susshi/internal/config"
)

func main() {
	if _, err := exec.LookPath("ssh"); err != nil {
		log.Fatalf("ssh not found: %v", err)
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	a, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}
	if err := a.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
