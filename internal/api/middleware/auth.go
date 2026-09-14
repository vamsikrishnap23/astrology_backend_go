package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwks *keyfunc.JWKS
var staticPubKey interface{}

// InitAuth initializes the JWKS fetcher or the static public key.
func InitAuth() error {
	// 1. Check if a static PEM public key was provided (Best for custom environments)
	pubKeyPEM := os.Getenv("SUPABASE_PUBLIC_KEY")
	if pubKeyPEM != "" {
		// Clean up escaped newlines if passed via certain environments
		pubKeyPEM = strings.ReplaceAll(pubKeyPEM, "\\n", "\n")

		// If the environment variable collapsed the newlines into spaces (common in Render/Docker),
		// we need to rebuild the valid PEM format.
		if !strings.Contains(pubKeyPEM, "\n") {
			pubKeyPEM = strings.ReplaceAll(pubKeyPEM, "-----BEGIN PUBLIC KEY-----", "-----BEGIN PUBLIC KEY-----\n")
			pubKeyPEM = strings.ReplaceAll(pubKeyPEM, "-----END PUBLIC KEY-----", "\n-----END PUBLIC KEY-----")

			// Replace any spaces inside the base64 payload with newlines
			// First, extract the middle part
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

	// 2. Fallback to fetching JWKS dynamically
	projectRef := os.Getenv("SUPABASE_PROJECT_REF")
	if projectRef == "" {
		log.Println("WARNING: Neither SUPABASE_PUBLIC_KEY nor SUPABASE_PROJECT_REF is set.")
		return nil
	}

	projectRef = strings.TrimPrefix(projectRef, "https://")
	projectRef = strings.TrimPrefix(projectRef, "http://")
	projectRef = strings.TrimSuffix(projectRef, ".supabase.co")
	projectRef = strings.TrimSuffix(projectRef, "/")

	jwksURL := "https://" + projectRef + ".supabase.co/auth/v1/jwks"

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

// SupabaseAuthMiddleware ensures the request has a valid Supabase JWT signed via ES256.
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

		// Token is valid! Proceed to the astrology calculations
		next.ServeHTTP(w, r)
	})
}
