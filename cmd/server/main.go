package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/vamsi/astrology_backend_go/internal/api/handlers"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/ephemeris"
)

func main() {
	// Initialize Ephemeris
	ephePath := os.Getenv("EPHE_PATH")
	if ephePath == "" {
		absPath, err := filepath.Abs("ephe_data")
		if err == nil {
			ephePath = absPath
		} else {
			ephePath = "ephe_data" // Fallback to relative if Abs fails
		}
	}

	// Swiss Ephemeris C library often strictly requires a trailing slash for directory paths
	if !strings.HasSuffix(ephePath, string(os.PathSeparator)) {
		ephePath += string(os.PathSeparator)
	}

	log.Printf("Resolved EPHE_PATH: %s", ephePath)

	// Force the environment variable so the C library can read it if swe_set_ephe_path fails
	os.Setenv("EPHE_PATH", ephePath)

	if err := ephemeris.Init(ephePath); err != nil {
		log.Fatalf("Failed to initialize ephemeris: %v", err)
	}
	defer ephemeris.Close()

	// Setup routing
	http.HandleFunc("/api/v1/chart", handlers.ChartHandler)
	http.HandleFunc("/api/v1/panchang", handlers.PanchangHandler)
	http.HandleFunc("/api/v1/tables", handlers.TablesHandler)
	http.HandleFunc("/api/v1/significators", handlers.SignificatorsHandler)
	http.HandleFunc("/api/v1/ruling-planets", handlers.RulingPlanetsHandler)
	http.HandleFunc("/api/v1/dasha", handlers.DashaHandler)
	http.HandleFunc("/api/v1/four-step", handlers.FourStepSignificatorsHandler)
	http.HandleFunc("/api/v1/vargas", handlers.VargasHandler)
	http.HandleFunc("/api/v1/progression", handlers.ProgressionHandler)
	http.HandleFunc("/api/v1/transits/chart", handlers.TransitHandler)
	http.HandleFunc("/api/v1/transits/upcoming", handlers.UpcomingTransitsHandler)
	http.HandleFunc("/api/v1/bhava-chalit", handlers.BhavaChalitHandler)
	http.HandleFunc("/api/v1/ashtakavarga", handlers.AshtakavargaHandler)
	http.HandleFunc("/api/v1/shadbala", handlers.ShadbalaHandler)
	http.HandleFunc("/api/v1/jaimini-karakas", handlers.JaiminiKarakasHandler)
	http.HandleFunc("/api/v1/ashtakoota", handlers.AshtakootaHandler)
	http.HandleFunc("/api/v1/btr", handlers.BTRHandler)

	// Serve static UI on root
	http.Handle("/", http.FileServer(http.Dir("static")))

	port := ":8080"
	log.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
