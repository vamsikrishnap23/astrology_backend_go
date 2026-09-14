package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwks *keyfunc.JWKS

// InitAuth initializes the JWKS fetcher. It should be called during server startup.
func InitAuth() error {
	projectRef := os.Getenv("SUPABASE_PROJECT_REF")
	if projectRef == "" {
		log.Println("WARNING: SUPABASE_PROJECT_REF not set. Auth middleware will fail if accessed.")
		return nil
	}

	jwksURL := "https://" + projectRef + ".supabase.co/auth/v1/jwks"

	// Create the JWKS from the resource at the given URL.
	options := keyfunc.Options{
		RefreshInterval: time.Hour,
		RefreshTimeout:  time.Second * 10,
		RefreshErrorHandler: func(err error) {
			log.Printf("There was an error with the jwt.Keyfunc\nError: %s", err.Error())
		},
	}

	var err error
	jwks, err = keyfunc.Get(jwksURL, options)
	if err != nil {
		log.Printf("Failed to create JWKS from resource at the given URL.\nError: %s", err.Error())
		return err
	}

	log.Printf("Successfully initialized Supabase JWKS from %s", jwksURL)
	return nil
}

// SupabaseAuthMiddleware ensures the request has a valid Supabase JWT signed via ES256 (or RS256).
func SupabaseAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow CORS preflight requests
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Missing Token", http.StatusUnauthorized)
			return
		}

		if jwks == nil {
			http.Error(w, "Server Configuration Error: JWKS not initialized", http.StatusInternalServerError)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and verify the token using the dynamically fetched JWKS
		token, err := jwt.Parse(tokenString, jwks.Keyfunc)
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid Token", http.StatusUnauthorized)
			return
		}

		// Token is valid! Proceed to the astrology calculations
		next.ServeHTTP(w, r)
	})
}
