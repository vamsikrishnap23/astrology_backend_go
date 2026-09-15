package middleware

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
)

var jwks *keyfunc.JWKS
var staticPubKey interface{}
var db *sql.DB

// InitAuth initializes the JWKS fetcher or the static public key, and the DB connection.
func InitAuth() error {
	// 1. Connect to Supabase DB for Admin Approvals
	dbUrl := os.Getenv("SUPABASE_DB_URL")
	if dbUrl != "" {
		var dbErr error
		db, dbErr = sql.Open("postgres", dbUrl)
		if dbErr != nil {
			log.Fatalf("Failed to open Supabase DB: %v", dbErr)
		}
		if err := db.Ping(); err != nil {
			log.Fatalf("Failed to ping Supabase DB: %v", err)
		}
		log.Println("Successfully connected to Supabase Database for Admin Approvals.")
	} else {
		log.Println("WARNING: SUPABASE_DB_URL is not set. Database approval check will fail if enforced.")
	}

	// 2. Check if a static PEM public key was provided
	pubKeyPEM := os.Getenv("SUPABASE_PUBLIC_KEY")
	if pubKeyPEM != "" {
		pubKeyPEM = strings.ReplaceAll(pubKeyPEM, "\\n", "\n")

		if !strings.Contains(pubKeyPEM, "\n") {
			pubKeyPEM = strings.ReplaceAll(pubKeyPEM, "-----BEGIN PUBLIC KEY-----", "-----BEGIN PUBLIC KEY-----\n")
			pubKeyPEM = strings.ReplaceAll(pubKeyPEM, "-----END PUBLIC KEY-----", "\n-----END PUBLIC KEY-----")

			parts := strings.Split(pubKeyPEM, "\n")
			if len(parts) >= 3 {
				header := parts[0]
				payload := parts[1]
				footer := parts[2]

				payload = strings.ReplaceAll(payload, " ", "\n")
				pubKeyPEM = header + "\n" + payload + "\n" + footer
			}
		}

		key, err := jwt.ParseECPublicKeyFromPEM([]byte(pubKeyPEM))
		if err != nil {
			log.Fatalf("Failed to parse SUPABASE_PUBLIC_KEY: %v", err)
		}
		staticPubKey = key
		log.Println("Successfully loaded static Supabase ES256 Public Key from environment.")
		return nil
	}

	// 3. Fallback to fetching JWKS dynamically
	projectRef := os.Getenv("SUPABASE_PROJECT_REF")
	if projectRef == "" {
		log.Println("WARNING: Neither SUPABASE_PUBLIC_KEY nor SUPABASE_PROJECT_REF is set.")
		return nil
	}

	projectRef = strings.TrimPrefix(projectRef, "https://")
	projectRef = strings.TrimPrefix(projectRef, "http://")
	projectRef = strings.TrimSuffix(projectRef, ".supabase.co")
	projectRef = strings.TrimSuffix(projectRef, "/")

	jwksURL := "https://" + projectRef + ".supabase.co/auth/v1/.well-known/jwks.json"

	anonKey := os.Getenv("SUPABASE_ANON_KEY")
	if anonKey == "" {
		log.Println("WARNING: SUPABASE_ANON_KEY not set for JWKS fetch.")
	}

	client := &http.Client{
		Transport: &headerTransport{
			Transport: http.DefaultTransport,
			Headers: map[string]string{
				"apikey": anonKey,
			},
		},
	}

	options := keyfunc.Options{
		Client:          client,
		RefreshInterval: time.Hour,
		RefreshTimeout:  time.Second * 10,
		RefreshErrorHandler: func(err error) {
			log.Printf("JWKS Refresh Error: %s", err.Error())
		},
	}

	var err error
	jwks, err = keyfunc.Get(jwksURL, options)
	if err != nil {
		log.Printf("Failed to fetch JWKS from %s. (Did Supabase disable this endpoint?)\nError: %s", jwksURL, err.Error())
		return err
	}

	log.Printf("Successfully initialized Supabase JWKS from %s", jwksURL)
	return nil
}

type headerTransport struct {
	Transport http.RoundTripper
	Headers   map[string]string
}

func (h *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.Headers {
		req.Header.Set(k, v)
	}
	return h.Transport.RoundTrip(req)
}

// SupabaseAuthMiddleware ensures the request has a valid Supabase JWT signed via ES256,
// and checks if the user's email is present in the approved_users table.
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

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		var keyFunc jwt.Keyfunc
		if staticPubKey != nil {
			keyFunc = func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return staticPubKey, nil
			}
		} else if jwks != nil {
			keyFunc = jwks.Keyfunc
		} else {
			http.Error(w, "Server Configuration Error: Auth keys not initialized", http.StatusInternalServerError)
			return
		}

		// Parse and verify the token
		token, err := jwt.Parse(tokenString, keyFunc)
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid Token", http.StatusUnauthorized)
			return
		}

		// Token is valid. Now check if the user's email is in the approved_users table.
		if db == nil {
			http.Error(w, "Server Configuration Error: Database not initialized", http.StatusInternalServerError)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized: Invalid Token Claims", http.StatusUnauthorized)
			return
		}

		email, ok := claims["email"].(string)
		if !ok || email == "" {
			http.Error(w, "Unauthorized: Email claim missing from token", http.StatusUnauthorized)
			return
		}

		var exists int
		queryErr := db.QueryRow("SELECT 1 FROM approved_users WHERE email = $1 LIMIT 1", email).Scan(&exists)
		if queryErr != nil {
			if queryErr == sql.ErrNoRows {
				http.Error(w, "Forbidden: User is not approved", http.StatusForbidden)
				return
			}
			log.Printf("Database error while checking approved_users: %v", queryErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Email exists in approved_users! Proceed to the astrology calculations
		next.ServeHTTP(w, r)
	})
}
