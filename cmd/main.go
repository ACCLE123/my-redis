package main

import (
	"log"
	"myredis/internal/server"
)

func main() {
	srv := server.New()
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}