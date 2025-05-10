package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"strathlearn/backend/auth"
	"strathlearn/backend/models"
	"strathlearn/backend/services/runner"
)

type APIHandler struct {
	challenges map[string]models.Challenge
	runner     runner.CodeRunner
}

func NewAPIHandler(challenges map[string]models.Challenge, runner runner.CodeRunner) *APIHandler {
	return &APIHandler{
		challenges: challenges,
		runner:     runner,
	}
}

func (h *APIHandler) ListChallenges(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.challenges)
}

func (h *APIHandler) GetChallenge(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/challenge/"):]
	log.Printf("Request for challenge: %s", id)
	if challenge, ok := h.challenges[id]; ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(challenge)
	} else {
		log.Printf("Challenge not found: %s", id)
		http.NotFound(w, r)
	}
}

func (h *APIHandler) SubmitSolution(w http.ResponseWriter, r *http.Request) {
	var req models.SubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	challenge, ok := h.challenges[req.ChallengeID]
	if !ok {
		http.NotFound(w, r)
		return
	}

	results := h.runner.RunTests(req.Code, challenge)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.SubmissionResponse{
		Success:     allTestsPassed(results),
		Message:     "Submission processed",
		TestResults: results,
	})
}

func (h *APIHandler) Debug(w http.ResponseWriter, r *http.Request) {
	currentDir, _ := os.Getwd()
	fmt.Fprintf(w, "Current directory: %s\n\n", currentDir)

	challengesDir := "./backend/challenges"
	fmt.Fprintf(w, "Contents of challenges directory:\n")
	files, err := os.ReadDir(challengesDir)
	if err != nil {
		fmt.Fprintf(w, "Error reading challenges dir: %v\n", err)
	} else {
		for _, file := range files {
			fmt.Fprintf(w, "- %s\n", file.Name())
		}
	}

	fmt.Fprintf(w, "\nLoaded challenges:\n")
	for id, challenge := range h.challenges {
		fmt.Fprintf(w, "- ID: %s, Title: %s, Source: %s\n",
			id, challenge.Title, challenge.FilePath)
	}

	fmt.Fprintf(w, "\nDocker connection status: %v\n", h.runner.IsDockerAvailable())
}

func allTestsPassed(results []models.TestResult) bool {
	for _, result := range results {
		if !result.Passed {
			return false
		}
	}
	return len(results) > 0
}

func (h *APIHandler) TestAuth(w http.ResponseWriter, r *http.Request) {

	user, ok := auth.GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Authentication successful",
		"user": map[string]interface{}{
			"id":            user.ID,
			"name":          user.Name,
			"email":         user.Email,
			"emailVerified": user.EmailVerified,
			"createdAt":     user.CreatedAt,
		},
	})
}
