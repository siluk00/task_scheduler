package main

import (
	"log"

	"github.com/siluk00/task_scheduler/internal/api"
	"github.com/siluk00/task_scheduler/pkg/config"
)

func main() {
	cfg := config.LoadConfig()

	//The server is returned with routes prepared
	server, err := api.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	log.Printf("Starting server on port %s...", cfg.ServerPort)
	// Begins the server in the port
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
