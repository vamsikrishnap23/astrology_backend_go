package main

import (
	"log"
	"net/http"
	"os"

	"github.com/vamsi/astrology_backend_go/internal/api/handlers"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/ephemeris"
)

func main() {
	// Initialize Ephemeris
	// By default, try to find in current directory or system path
	ephePath := os.Getenv("EPHE_PATH")
	ephemeris.Init(ephePath)
	defer ephemeris.Close()

	http.HandleFunc("/api/v1/chart", handlers.ChartHandler)

	port := ":8080"
	log.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
