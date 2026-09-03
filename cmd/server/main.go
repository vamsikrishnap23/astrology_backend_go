package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/api/handlers"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
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

	// Health Check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Setup routing
	http.HandleFunc("/api/chart", handlers.ChartHandler)
	http.HandleFunc("/api/panchang", handlers.PanchangHandler)
	http.HandleFunc("/api/tables", handlers.TablesHandler)
	http.HandleFunc("/api/significators", handlers.SignificatorsHandler)
	http.HandleFunc("/api/ruling-planets", handlers.RulingPlanetsHandler)
	http.HandleFunc("/api/dasha", handlers.DashaHandler)
	http.HandleFunc("/api/four-step", handlers.FourStepSignificatorsHandler)
	http.HandleFunc("/api/vargas", handlers.VargasHandler)
	http.HandleFunc("/api/progression", handlers.ProgressionHandler)
	http.HandleFunc("/api/transits/chart", handlers.TransitHandler)
	http.HandleFunc("/api/transits/upcoming", handlers.UpcomingTransitsHandler)
	http.HandleFunc("/api/bhava-chalit", handlers.BhavaChalitHandler)
	http.HandleFunc("/api/ashtakavarga", handlers.AshtakavargaHandler)
	http.HandleFunc("/api/shadbala", handlers.ShadbalaHandler)
	http.HandleFunc("/api/jaimini-karakas", handlers.JaiminiKarakasHandler)
	http.HandleFunc("/api/ashtakoota", handlers.AshtakootaHandler)
	http.HandleFunc("/api/btr", handlers.BTRHandler)

	// Serve static UI on root
	http.Handle("/", http.FileServer(http.Dir("static")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + strings.TrimPrefix(port, ":")
	log.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
