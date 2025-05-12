package auth

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWKS represents the JSON Web Key Set structure
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWK represents a JSON Web Key
type JWK struct {
	Crv string `json:"crv"`
	X   string `json:"x"`
	Kty string `json:"kty"`
	Kid string `json:"kid"`
}

// JWTClaims represents the claims in our JWT
type JWTClaims struct {
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"emailVerified"`
	Image         *string   `json:"image"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ID            string    `json:"id"`
	jwt.RegisteredClaims
}

// JWKSCache stores the JWKS and manages its expiration
type JWKSCache struct {
	jwks          *JWKS
	fetchedAt     time.Time
	cacheDuration time.Duration
	mu            sync.RWMutex
}

// Global JWKS cache with 24-hour expiration
var jwksCache = NewJWKSCache(24 * time.Hour)

// NewJWKSCache creates a new JWKS cache with the given cache duration
func NewJWKSCache(cacheDuration time.Duration) *JWKSCache {
	return &JWKSCache{
		cacheDuration: cacheDuration,
	}
}

// GetJWKS gets the JWKS, fetching it if necessary
func (c *JWKSCache) GetJWKS(jwksURL string) (*JWKS, error) {
	c.mu.RLock()
	if c.jwks != nil && time.Since(c.fetchedAt) < c.cacheDuration {
		jwks := c.jwks
		c.mu.RUnlock()
		return jwks, nil
	}
	c.mu.RUnlock()

	// Need to refresh
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check again in case another goroutine updated while we were waiting
	if c.jwks != nil && time.Since(c.fetchedAt) < c.cacheDuration {
		return c.jwks, nil
	}

	// Fetch new JWKS
	jwks, err := fetchJWKS(jwksURL)
	if err != nil {
		return nil, err
	}

	c.jwks = jwks
	c.fetchedAt = time.Now()
	return jwks, nil
}

// fetchJWKS fetches the JWKS from the SvelteKit server
func fetchJWKS(jwksURL string) (*JWKS, error) {
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch JWKS: status code %d", resp.StatusCode)
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %v", err)
	}

	return &jwks, nil
}

// AuthMiddleware verifies JWT tokens
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		// Get the token
		tokenString := parts[1]

		// Parse the token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Get the key ID from the token header
			kidInterface, ok := token.Header["kid"]
			if !ok {
				return nil, fmt.Errorf("kid not found in token header")
			}
			kid, ok := kidInterface.(string)
			if !ok {
				return nil, fmt.Errorf("kid is not a string")
			}

			// Get JWKS from cache or fetch from SvelteKit server
			jwksURL := "https://codex.singularity.co.ke/api/auth/jwks"
			jwks, err := jwksCache.GetJWKS(jwksURL)
			if err != nil {
				return nil, err
			}

			// Find the key with matching kid
			var key *JWK
			for _, k := range jwks.Keys {
				if k.Kid == kid {
					key = &k
					break
				}
			}
			if key == nil {
				return nil, fmt.Errorf("no matching key found for kid: %s", kid)
			}

			// Decode the public key
			if key.Kty != "OKP" || key.Crv != "Ed25519" {
				return nil, fmt.Errorf("unsupported key type or curve")
			}

			xBytes, err := base64.RawURLEncoding.DecodeString(key.X)
			if err != nil {
				return nil, fmt.Errorf("failed to decode public key: %v", err)
			}

			// Ed25519 public key is 32 bytes
			if len(xBytes) != ed25519.PublicKeySize {
				return nil, fmt.Errorf("invalid public key size")
			}

			return ed25519.PublicKey(xBytes), nil
		})

		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Verify token is valid
		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims for use in handlers
		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
			return
		}

		// Add user info to request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", claims)
		r = r.WithContext(ctx)

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}

// GetUserFromContext extracts the user claims from the request context
func GetUserFromContext(r *http.Request) (*JWTClaims, bool) {
	user, ok := r.Context().Value("user").(*JWTClaims)
	return user, ok
}
