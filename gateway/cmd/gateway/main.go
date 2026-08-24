package main

import (
	"gateway/internal/config"
	"gateway/internal/server"
	"log"
)

func main() {
	cfg, err := config.Load("/config/config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := server.ListenAndServe(cfg); err != nil {
		log.Fatal(err)
	}
}
