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
	ephePath := os.Getenv("EPHE_PATH")
	if ephePath == "" {
		ephePath = "ephe_data" // Fallback to local directory
	}
	if err := ephemeris.Init(ephePath); err != nil {
		log.Fatalf("Failed to initialize ephemeris: %v", err)
	}
	defer ephemeris.Close()

	// Setup routing
	http.HandleFunc("/api/v1/chart", handlers.ChartHandler)
	http.HandleFunc("/api/v1/panchang", handlers.PanchangHandler)
	http.HandleFunc("/api/v1/tables", handlers.TablesHandler)

	port := ":8080"
	log.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
