package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"strathlearn/backend/auth"
	"strathlearn/backend/handlers"
	"strathlearn/backend/services"
	"strathlearn/backend/services/runner"

	"github.com/docker/docker/client"
)

var dockerClient *client.Client

func main() {
	log.Println("Starting server...")
	currentDir, _ := os.Getwd()
	log.Printf("Current working directory: %s", currentDir)

	var err error
	dockerClient, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Printf("Warning: Could not connect to Docker: %v", err)
		log.Println("Falling back to local execution mode")
	} else {
		log.Println("Successfully connected to Docker")
		go runner.StartContainerCleanupWorker(dockerClient)
	}

	challengesDir := "./backend/challenges"
	if _, err := os.Stat(challengesDir); os.IsNotExist(err) {
		log.Printf("Creating challenges directory at %s", challengesDir)
		tempDir, err := os.MkdirTemp(".", "challenges")
		if err != nil {
			log.Fatalf("Failed to create challenges directory: %v", err)
		}
		challengesDir = tempDir
		services.CreateSampleChallenge(challengesDir)
	}

	challenges := services.LoadChallenges(challengesDir)
	log.Printf("Loaded %d challenges", len(challenges))

	for id, challenge := range challenges {
		log.Printf("Challenge in memory: ID=%s, Title=%s, Source=%s",
			id, challenge.Title, challenge.FilePath)
	}

	codeRunner := runner.NewRunner(dockerClient)
	apiHandler := handlers.NewAPIHandler(challenges, codeRunner)

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	http.HandleFunc("/api/challenges", apiHandler.ListChallenges)

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/challenge/", apiHandler.GetChallenge)
	protectedMux.HandleFunc("/api/submit", apiHandler.SubmitSolution)
	protectedMux.HandleFunc("/debug", apiHandler.Debug)
	protectedMux.HandleFunc("/api/test-auth", apiHandler.TestAuth)

	authHandler := corsMiddleware(auth.AuthMiddleware(protectedMux))
	http.Handle("/api/challenge/", authHandler)
	http.Handle("/api/submit", authHandler)
	http.Handle("/debug", authHandler)

	http.Handle("/api/test-auth", authHandler)

	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
