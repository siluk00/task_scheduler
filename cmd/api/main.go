package main

import (
	"log"

	"github.com/siluk00/task_scheduler/internal/api"
	"github.com/siluk00/task_scheduler/pkg/config"
)

func main() {
	cfg := config.LoadConfig()

	//O servidor é retornado com as rotas configuradas e pronto para ser iniciado.
	server, err := api.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	log.Printf("Starting server on port %s...", cfg.ServerPort)
	// O método Start inicia o servidor HTTP na porta configurada.
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
