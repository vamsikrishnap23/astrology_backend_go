// Package site holds the public endpoints used by the company website:
// the daily panchang, the D1 (Rasi) chart and the list of every value they
// can return (with Telugu translations).
//
// Unlike /api/*, these routes need no login. Settings are fixed: Lahiri
// ayanamsa for both, Placidus houses for the chart.
package site

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

// Register adds the website routes to mux:
//
//	GET|POST /site/panchang   daily panchang
//	POST     /site/chart      D1 chart
//	GET      /site/values     every possible value + Telugu translation
func Register(mux *http.ServeMux) {
	origins := allowedOrigins()
	mux.Handle("/site/panchang", withCORS(origins, http.HandlerFunc(PanchangHandler)))
	mux.Handle("/site/chart", withCORS(origins, http.HandlerFunc(ChartHandler)))
	mux.Handle("/site/values", withCORS(origins, http.HandlerFunc(ValuesHandler)))
}

// allowedOrigins reads SITE_ALLOWED_ORIGINS (comma separated, e.g.
// "https://example.com,https://www.example.com"). Empty means any origin.
func allowedOrigins() []string {
	var out []string
	for _, o := range strings.Split(os.Getenv("SITE_ALLOWED_ORIGINS"), ",") {
		if o = strings.TrimSpace(o); o != "" {
			out = append(out, strings.TrimSuffix(o, "/"))
		}
	}
	return out
}

func withCORS(origins []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if len(origins) == 0 {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else if origin != "" {
			for _, o := range origins {
				if o == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Add("Vary", "Origin")
					break
				}
			}
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// writeJSON sends v as JSON with the given status.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError sends {"error": msg}.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// validLocation checks latitude, longitude and timezone ranges.
func validLocation(lat, lon, tz float64) string {
	switch {
	case lat < -90 || lat > 90:
		return "latitude must be between -90 and 90"
	case lon < -180 || lon > 180:
		return "longitude must be between -180 and 180"
	case tz < -14 || tz > 14:
		return "timezone must be a UTC offset in hours between -14 and 14 (e.g. 5.5)"
	}
	return ""
}

// validYear keeps requests inside the bundled ephemeris files (1800-2399).
func validYear(y int) string {
	if y < 1800 || y > 2399 {
		return "date must be between 1800 and 2399"
	}
	return ""
}
