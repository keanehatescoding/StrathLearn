package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// User represents an authenticated user
type User struct {
	ID       string    `json:"id"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Provider string    `json:"provider"`
	LastSeen time.Time `json:"lastSeen"`
}

// Auth handles authentication
type Auth struct {
	Enabled     bool
	ProviderURL string
	CallbackURL string
	JWTSecret   string
	Users       map[string]User // In-memory user store (replace with DB in production)
}

// NewAuth creates a new Auth service
func NewAuth(enabled bool, providerURL, callbackURL, jwtSecret string) *Auth {
	return &Auth{
		Enabled:     enabled,
		ProviderURL: providerURL,
		CallbackURL: callbackURL,
		JWTSecret:   jwtSecret,
		Users:       make(map[string]User),
	}
}

// Middleware is authentication middleware that checks for a valid JWT
func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for public paths
		if !a.Enabled || isPublicPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(a.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Token is valid, continue to handler
		next.ServeHTTP(w, r)
	})
}

// isPublicPath determines if a path is publicly accessible without auth
func isPublicPath(path string) bool {
	publicPaths := []string{
		"/",
		"/index.html",
		"/api/auth/",
		"/favicon.ico",
	}

	// Check if path starts with any public path prefix
	for _, prefix := range publicPaths {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	// Also allow static files
	if strings.HasPrefix(path, "/static/") ||
		strings.HasPrefix(path, "/assets/") ||
		strings.HasPrefix(path, "/css/") ||
		strings.HasPrefix(path, "/js/") {
		return true
	}

	return false
}

// GetUser retrieves the current user from a request's JWT
func (a *Auth) GetUser(r *http.Request) (*User, error) {
	if !a.Enabled {
		return &User{ID: "anonymous"}, nil
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("no token provided")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	user, exists := a.Users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

// RegisterAuthRoutes sets up authentication endpoints
func (a *Auth) RegisterAuthRoutes(mux *http.ServeMux) {
	if !a.Enabled {
		return
	}

	// This would handle login, callback from auth provider, etc.
	// For Better Auth integration, you'd implement their specific flows here
	mux.HandleFunc("/api/auth/login", a.handleLogin)
	mux.HandleFunc("/api/auth/callback", a.handleCallback)
	mux.HandleFunc("/api/auth/user", a.handleGetUser)
}

// handleLogin initiates the authentication flow
func (a *Auth) handleLogin(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, you would redirect to Better Auth login page
	// For now we'll just return info about the auth config
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"provider": a.ProviderURL,
		"callback": a.CallbackURL,
		"message":  "Auth integration would start here",
	})
}

// handleCallback processes the callback from auth provider
func (a *Auth) handleCallback(w http.ResponseWriter, r *http.Request) {
	// This would process the callback from Better Auth
	// For now, we'll create a dummy user and token

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing auth code", http.StatusBadRequest)
		return
	}

	// In a real implementation, you would exchange the code for tokens
	// with Better Auth and extract user info

	userID := "user_" + code[:8]
	user := User{
		ID:       userID,
		Email:    userID + "@example.com",
		Name:     "Test User",
		Provider: "better-auth",
		LastSeen: time.Now(),
	}

	a.Users[userID] = user

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userID,
		"name":  user.Name,
		"email": user.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(a.JWTSecret))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return token to frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
		"user":  user.Name,
	})
}

// handleGetUser returns the current user's info
func (a *Auth) handleGetUser(w http.ResponseWriter, r *http.Request) {
	user, err := a.GetUser(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
