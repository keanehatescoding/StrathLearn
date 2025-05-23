package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"strathlearn/backend/auth"
	"strathlearn/backend/data" // Provides data.GetLanguages
	"strathlearn/backend/db"
	"strathlearn/backend/handlers"
	"strathlearn/backend/services"
	"strathlearn/backend/services/runner" // Provides runner.ExecuteCode & NewJudge0Runner

	"github.com/docker/docker/client" // If still used by other parts
	"github.com/gin-gonic/gin"
)

var dockerClient *client.Client

func main() {
	log.Println("Starting server...")
	log.Println("Connecting to database...")
	dbConn, dbErr := db.Connect()
	if dbErr != nil {
		log.Fatalf("Failed to connect to database: %v", dbErr)
	}
	defer db.Disconnect(dbConn)
	currentDir, _ := os.Getwd()
	log.Printf("Current working directory: %s", currentDir)

	// Initialize Docker client if needed by other parts of your application
	// var err error
	// dockerClient, err = client.NewClientWithOpts(client.FromEnv)
	// if err != nil {
	// 	log.Printf("Warning: Could not connect to Docker: %v", err)
	// } else {
	// 	log.Println("Successfully connected to Docker")
	// 	// go runner.StartContainerCleanupWorker(dockerClient) // If you have such a worker
	// }

	challengesDir := "./backend/challenges"
	if _, err := os.Stat(challengesDir); os.IsNotExist(err) {
		log.Printf("Creating challenges directory at %s", challengesDir)
		// Simplified: ensure challengesDir exists or handle appropriately
		// tempDir, err := os.MkdirTemp(".", "challenges")
		// if err != nil {
		// 	log.Fatalf("Failed to create challenges directory: %v", err)
		// }
		// challengesDir = tempDir
		// services.CreateSampleChallenge(challengesDir)
	}

	challenges := services.LoadChallenges(challengesDir)
	log.Printf("Loaded %d challenges", len(challenges))

	for id, challenge := range challenges {
		log.Printf("Challenge in memory: ID=%s, Title=%s, Source=%s",
			id, challenge.Title, challenge.FilePath)
	}

	// Initialize your existing API handler and runner
	// The runner for apiHandler should be the Judge0 runner instance.
	judge0Runner := runner.NewJudge0Runner()
	apiHandler := handlers.NewAPIHandler(challenges, judge0Runner) // Pass the correct runner

	// --- CORS Middleware (Your existing middleware) ---
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedOrigins := []string{
				"https://codex.singularity.co.ke",
				"http://localhost:5173",
				"http://localhost:3000",
				"http://localhost:8080",
			}
			origin := r.Header.Get("Origin")
			isAllowed := false
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					isAllowed = true
					break
				}
			}

			if isAllowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if origin == "" {
				w.Header().Set("Access-Control-Allow-Origin", "*") // Or handle more strictly
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	// --- Initialize Gin Engine for new Gin-based handlers ---
	// gin.Default() includes Logger and Recovery middleware.
	// For new public routes that use Gin handlers:
	publicGinRouter := gin.Default()
	// For new protected routes that use Gin handlers:
	protectedGinRouter := gin.Default()

	// --- Register NEW PUBLIC Gin handlers ---
	// Example: GET /api/languages
	// The Gin router will handle the exact path "/api/languages"
	publicGinRouter.GET("/api/languages", data.GetLanguages)
	// Since publicGinRouter is an http.Handler, we can use it with http.Handle
	// The entire publicGinRouter will handle requests matching this path prefix.
	// Gin's internal router will then dispatch to the correct handler based on method and path.
	http.Handle("/api/languages", corsMiddleware(publicGinRouter))

	// --- Setup for EXISTING and NEW PROTECTED routes ---
	// Your existing ServeMux for protected routes
	protectedMux := http.NewServeMux()

	// Add EXISTING protected http.HandlerFunc routes to protectedMux
	// (These are assumed to be http.HandlerFunc compatible as per your original code structure)
	protectedMux.HandleFunc("/api/challenge/", apiHandler.GetChallenge)                // Example, adjust as needed
	protectedMux.HandleFunc("/api/submit", apiHandler.SubmitSolution)                  // Example
	protectedMux.HandleFunc("/debug", apiHandler.Debug)                                // Example
	protectedMux.HandleFunc("/api/test-auth", apiHandler.TestAuth)                     // Example
	protectedMux.HandleFunc("/api/profile", apiHandler.GetUserProfile)                 // Example
	protectedMux.HandleFunc("/api/profile/submissions", apiHandler.GetUserSubmissions) //Example

	// --- Register NEW PROTECTED Gin handlers ---
	// Example: POST /api/execute-code
	// The Gin router will handle the exact path "/api/execute-code"
	protectedGinRouter.POST("/api/execute-code", runner.ExecuteCode)
	// Add the protectedGinRouter as a handler in your existing protectedMux
	// When a request hits "/api/execute-code" on protectedMux, it will be passed to protectedGinRouter.
	protectedMux.Handle("/api/execute-code", protectedGinRouter) // protectedGinRouter is an http.Handler

	// --- Apply Auth Middleware to all routes in protectedMux ---
	authHandler := auth.AuthMiddleware(protectedMux) // auth.AuthMiddleware wraps protectedMux

	// --- Apply CORS and then Auth to protected route groups/paths ---
	// These Handle calls link paths to the (cors(auth(protectedMux))) chain.
	// Ensure these paths correctly cover all routes handled by protectedMux.
	http.Handle("/api/challenge/", corsMiddleware(authHandler))
	http.Handle("/api/submit", corsMiddleware(authHandler))
	http.Handle("/debug", corsMiddleware(authHandler))
	http.Handle("/api/test-auth", corsMiddleware(authHandler))
	http.Handle("/api/profile", corsMiddleware(authHandler))
	http.Handle("/api/profile/submissions", corsMiddleware(authHandler))
	// Add a handler for the new protected base path if not covered by broader patterns
	http.Handle("/api/execute-code", corsMiddleware(authHandler)) // Ensures this path goes through the chain

	// --- Handle other PUBLIC routes (existing, non-Gin) ---
	http.Handle("/api/challenges", corsMiddleware(http.HandlerFunc(apiHandler.ListChallenges)))

	// --- Static file server (existing) ---
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Printf("Server running on port %s\n", port)
	// http.ListenAndServe uses http.DefaultServeMux if the second argument is nil
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
