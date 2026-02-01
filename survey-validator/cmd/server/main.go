package main

import (
	"flag"
	"log"
	"os"

	"github.com/survey-validator/api"
)

func main() {
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	if envPort := os.Getenv("PORT"); envPort != "" {
		*port = envPort
	}

	addr := ":" + *port
	server := api.NewServer(addr)

	log.Println("Survey Data Validation & Insight Engine")
	log.Println("========================================")
	log.Printf("Server starting on http://localhost%s", addr)
	log.Println("Endpoints:")
	log.Println("  GET  /health           - Health check")
	log.Println("  POST /api/v1/validate  - Validate survey data")
	log.Println("========================================")

	if err := server.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
