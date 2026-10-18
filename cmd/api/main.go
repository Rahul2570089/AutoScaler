package main

import (
	"AutoScaler/config"
	"AutoScaler/internal/api"
	"AutoScaler/internal/storage"
	"fmt"
	"log"
)

func main() {
	cfg := config.DefaultConfig()

	store, err := storage.NewBadgerStorage(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}
	defer store.Close()

	server := api.NewServer(store)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting API server on %s", addr)

	if err := server.Run(addr); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
