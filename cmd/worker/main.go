package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/siluk00/task_scheduler/internal/worker"
	"github.com/siluk00/task_scheduler/pkg/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.LoadConfig()

	//Initialize worker
	taskWorker, err := worker.NewTaskWorker(cfg)
	if err != nil {
		log.Fatalf("Worker had poblems: %v", err.Error())
	}

	//Waiting for CTRL+C to finish
	sigChan := make(chan os.Signal, 1)
	//signal.Notify registers the given channel to receive notifications of specified signals.
	// syscall.SIGINT and syscall.SIGTERM are signals that indicate an interrupt or termination request.
	// SIGINT is typically sent when the user presses Ctrl+C in the terminal,
	// while SIGTERM is a termination signal sent by the operating system.
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := taskWorker.Start(ctx); err != nil {
			log.Printf("Worker stopped with error: %v", err.Error())
			sigChan <- syscall.SIGTERM
		}
	}()

	<-sigChan
	log.Println("Shutting down worker...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	taskWorker.Stop(shutdownCtx)
	log.Println("Worker stopped gracefully")
}
