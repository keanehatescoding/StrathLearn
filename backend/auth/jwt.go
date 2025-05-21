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

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Crv string `json:"crv"`
	X   string `json:"x"`
	Kty string `json:"kty"`
	Kid string `json:"kid"`
}

type JWTClaims struct {
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	EmailVerified bool       `json:"emailVerified"`
	Image         *string    `json:"image"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	Role          string     `json:"role"`
	Banned        *bool      `json:"banned"`
	BanReason     *string    `json:"banReason"`
	BanExpires    *time.Time `json:"banExpires"`
	ID            string     `json:"id"`
	jwt.RegisteredClaims
}

type JWKSCache struct {
	jwks          *JWKS
	fetchedAt     time.Time
	cacheDuration time.Duration
	mu            sync.RWMutex
}

var jwksCache = NewJWKSCache(24 * time.Hour)

func NewJWKSCache(cacheDuration time.Duration) *JWKSCache {
	return &JWKSCache{
		cacheDuration: cacheDuration,
	}
}

func (c *JWKSCache) GetJWKS(jwksURL string) (*JWKS, error) {
	c.mu.RLock()
	if c.jwks != nil && time.Since(c.fetchedAt) < c.cacheDuration {
		jwks := c.jwks
		c.mu.RUnlock()
		return jwks, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.jwks != nil && time.Since(c.fetchedAt) < c.cacheDuration {
		return c.jwks, nil
	}

	jwks, err := fetchJWKS(jwksURL)
	if err != nil {
		return nil, err
	}

	c.jwks = jwks
	c.fetchedAt = time.Now()
	return jwks, nil
}

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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			kidInterface, ok := token.Header["kid"]
			if !ok {
				return nil, fmt.Errorf("kid not found in token header")
			}
			kid, ok := kidInterface.(string)
			if !ok {
				return nil, fmt.Errorf("kid is not a string")
			}

			jwksURL := "https://codex.singularity.co.ke/api/auth/jwks"
			jwks, err := jwksCache.GetJWKS(jwksURL)
			if err != nil {
				return nil, err
			}

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

			if key.Kty != "OKP" || key.Crv != "Ed25519" {
				return nil, fmt.Errorf("unsupported key type or curve")
			}

			xBytes, err := base64.RawURLEncoding.DecodeString(key.X)
			if err != nil {
				return nil, fmt.Errorf("failed to decode public key: %v", err)
			}

			if len(xBytes) != ed25519.PublicKeySize {
				return nil, fmt.Errorf("invalid public key size")
			}

			return ed25519.PublicKey(xBytes), nil
		})

		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func GetUserFromContext(r *http.Request) (*JWTClaims, bool) {
	user, ok := r.Context().Value("user").(*JWTClaims)
	return user, ok
}
