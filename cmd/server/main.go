package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/api/handlers"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/api/middleware"
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

	// Setup routing for API with Auth Middleware
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/chart", handlers.ChartHandler)
	apiMux.HandleFunc("/api/panchang", handlers.PanchangHandler)
	apiMux.HandleFunc("/api/panchang/daily", handlers.DailyPanchangHandler)
	apiMux.HandleFunc("/api/tables", handlers.TablesHandler)
	apiMux.HandleFunc("/api/significators", handlers.SignificatorsHandler)
	apiMux.HandleFunc("/api/ruling-planets", handlers.RulingPlanetsHandler)
	apiMux.HandleFunc("/api/dasha", handlers.DashaHandler)
	apiMux.HandleFunc("/api/four-step", handlers.FourStepSignificatorsHandler)
	apiMux.HandleFunc("/api/vargas", handlers.VargasHandler)
	apiMux.HandleFunc("/api/progression", handlers.ProgressionHandler)
	apiMux.HandleFunc("/api/transits/chart", handlers.TransitHandler)
	apiMux.HandleFunc("/api/transits/rasi", handlers.RasiTransitHandler)
	apiMux.HandleFunc("/api/transits/upcoming", handlers.UpcomingTransitsHandler)
	apiMux.HandleFunc("/api/bhava-chalit", handlers.BhavaChalitHandler)
	apiMux.HandleFunc("/api/ashtakavarga", handlers.AshtakavargaHandler)
	apiMux.HandleFunc("/api/shadbala", handlers.ShadbalaHandler)
	apiMux.HandleFunc("/api/jaimini-karakas", handlers.JaiminiKarakasHandler)
	apiMux.HandleFunc("/api/ashtakoota", handlers.AshtakootaHandler)
	apiMux.HandleFunc("/api/btr", handlers.BTRHandler)
	apiMux.HandleFunc("/api/retrogrades", handlers.RetrogradesHandler)
	apiMux.HandleFunc("/api/manglik-dosha", handlers.ManglikHandler)

	// Wrap apiMux with SupabaseAuthMiddleware
	http.Handle("/api/", middleware.SupabaseAuthMiddleware(apiMux))

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
