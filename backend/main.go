package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Challenge struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Difficulty  string     `json:"difficulty"`
	Description string     `json:"description"`
	Hints       []string   `json:"hints"`
	TestCases   []TestCase `json:"testCases"`
	InitialCode string     `json:"initialCode"`
	Solutions   []string   `json:"solutions"`
	TimeLimit   int        `json:"timeLimit"`
	MemoryLimit int        `json:"memoryLimit"`
	FilePath    string     `json:"-"` // Track source file path for debugging
}

type TestCase struct {
	ID             string `json:"id"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expectedOutput"`
	Hidden         bool   `json:"hidden"`
}

type SubmissionRequest struct {
	ChallengeID string `json:"challengeId"`
	Code        string `json:"code"`
}

type TestResult struct {
	TestCaseID string `json:"testCaseId"`
	Passed     bool   `json:"passed"`
	Output     string `json:"output"`
	Error      string `json:"error"`
}

type SubmissionResponse struct {
	Success     bool         `json:"success"`
	Message     string       `json:"message"`
	TestResults []TestResult `json:"testResults"`
}

var challenges map[string]Challenge

func main() {
	log.Println("Starting server...")
	currentDir, _ := os.Getwd()
	log.Printf("Current working directory: %s", currentDir)

	// Create challenges directory if it doesn't exist
	challengesDir := "./challenges"
	if _, err := os.Stat(challengesDir); os.IsNotExist(err) {
		log.Printf("Creating challenges directory at %s", challengesDir)
		if err := os.Mkdir(challengesDir, 0755); err != nil {
			log.Fatalf("Failed to create challenges directory: %v", err)
		}
		createSampleChallenge(challengesDir)
	}

	// Load challenges
	challenges = loadChallenges(challengesDir)
	log.Printf("Loaded %d challenges", len(challenges))

	// Print loaded challenges for debugging
	for id, challenge := range challenges {
		log.Printf("Challenge in memory: ID=%s, Title=%s, Source=%s",
			id, challenge.Title, challenge.FilePath)
	}

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	http.HandleFunc("/api/challenges", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(challenges)
	})

	http.HandleFunc("/api/challenge/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/api/challenge/"):]
		log.Printf("Request for challenge: %s", id)
		if challenge, ok := challenges[id]; ok {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(challenge)
		} else {
			log.Printf("Challenge not found: %s", id)
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/api/submit", func(w http.ResponseWriter, r *http.Request) {
		var req SubmissionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		challenge, ok := challenges[req.ChallengeID]
		if !ok {
			http.NotFound(w, r)
			return
		}

		results := runTests(req.Code, challenge)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SubmissionResponse{
			Success:     allTestsPassed(results),
			Message:     "Submission processed",
			TestResults: results,
		})
	})

	// Add debug endpoint
	http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Current directory: %s\n\n", currentDir)

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
		for id, challenge := range challenges {
			fmt.Fprintf(w, "- ID: %s, Title: %s, Source: %s\n",
				id, challenge.Title, challenge.FilePath)
		}
	})

	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsMiddleware(http.DefaultServeMux)))
}

func createSampleChallenge(dir string) {
	challengeJSON := `{
        "id": "hello-world",
        "title": "Hello, Lol",
        "difficulty": "beginner",
        "description": "Welcome to your first C programming challenge! Write a simple C program that prints the message 'Hello, World!' to the console.\n\nThis is the traditional first program for beginners in any programming language, and it will help you verify that your development environment is set up correctly.",
        "hints": [
            "Use the printf function from the stdio.h library to output text",
            "Don't forget to include the stdio.h header at the top of your program",
            "Remember that your main function should return an integer (typically 0 for successful execution)",
            "In C, strings need to be enclosed in double quotes"
        ],
        "testCases": [
            {
                "id": "test1",
                "input": "",
                "expectedOutput": "Hello, World!\\n",
                "hidden": false
            }
        ],
        "initialCode": "#include <stdio.h>\\n\\nint main() {\\n    // Write your code here\\n    \\n    return 0;\\n}",
        "solutions": [
            "#include <stdio.h>\\n\\nint main() {\\n    printf(\\\"Hello, World!\\\\n\\\");\\n    return 0;\\n}"
        ],
        "timeLimit": 1,
        "memoryLimit": 128
    }`

	filePath := filepath.Join(dir, "hello-world.json")
	err := os.WriteFile(filePath, []byte(challengeJSON), 0644)
	if err != nil {
		log.Fatalf("Failed to write sample challenge: %v", err)
	}
	log.Printf("Created sample challenge at: %s", filePath)
}

func loadChallenges(dir string) map[string]Challenge {
	challenges := make(map[string]Challenge)
	log.Printf("Loading challenges from directory: %s", dir)

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("Error reading challenges directory: %v", err)
		return challenges
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(dir, file.Name())
			log.Printf("Processing file: %s", filePath)

			data, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Error reading file %s: %v", filePath, err)
				continue
			}

			// Log a preview of the file content
			preview := string(data)
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			log.Printf("File content preview: %s", preview)

			var challenge Challenge
			if err := json.Unmarshal(data, &challenge); err != nil {
				log.Printf("JSON parse error in %s: %v", filePath, err)
				continue
			}

			// Store file path for debugging
			challenge.FilePath = filePath

			// Ensure ID exists
			if challenge.ID == "" {
				log.Printf("Warning: Challenge in %s has no ID, using filename", filePath)
				challenge.ID = strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			}

			log.Printf("Loaded challenge from %s: ID=%s, Title=%s",
				filePath, challenge.ID, challenge.Title)

			// Store challenge by its ID from the file
			challenges[challenge.ID] = challenge
		}
	}

	// Create sample challenge if no challenges found
	if len(challenges) == 0 {
		log.Println("No challenges found, creating sample challenge")
		createSampleChallenge(dir)
		// Recursive call to load the newly created challenge
		return loadChallenges(dir)
	}

	return challenges
}

func runTests(code string, challenge Challenge) []TestResult {
	results := make([]TestResult, 0, len(challenge.TestCases))
	log.Printf("Running tests for challenge: %s", challenge.ID)

	tempDir, err := os.MkdirTemp("", "challenge-")
	if err != nil {
		log.Printf("Failed to create temp dir: %v", err)
		return []TestResult{{
			TestCaseID: "compile",
			Passed:     false,
			Error:      "System error: could not create temporary directory",
		}}
	}
	defer os.RemoveAll(tempDir)

	submissionID := uuid.New().String()

	sourcePath := filepath.Join(tempDir, fmt.Sprintf("solution-%s.c", submissionID))
	if err := os.WriteFile(sourcePath, []byte(code), 0644); err != nil {
		return []TestResult{{
			TestCaseID: "compile",
			Passed:     false,
			Error:      "System error: could not write source code",
		}}
	}

	execPath := filepath.Join(tempDir, fmt.Sprintf("solution-%s", submissionID))
	cmd := exec.Command("gcc", sourcePath, "-o", execPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return []TestResult{{
			TestCaseID: "compile",
			Passed:     false,
			Error:      "Compilation error: " + string(output),
		}}
	}

	for _, tc := range challenge.TestCases {
		result := TestResult{
			TestCaseID: tc.ID,
			Passed:     false,
		}

		cmd := exec.Command(execPath)
		cmd.Stdin = strings.NewReader(tc.Input)

		// Set a timeout
		timeout := time.Duration(challenge.TimeLimit) * time.Second
		if timeout == 0 {
			timeout = 5 * time.Second // Default timeout
		}

		outputChan := make(chan []byte, 1)
		errorChan := make(chan error, 1)

		go func() {
			output, err := cmd.CombinedOutput()
			outputChan <- output
			errorChan <- err
		}()

		// Timeout handling
		select {
		case output := <-outputChan:
			err := <-errorChan
			if err != nil {
				result.Error = "Runtime error: " + err.Error()
			} else {
				result.Output = string(output)
				outputTrimmed := strings.TrimSpace(result.Output)
				expectedTrimmed := strings.TrimSpace(tc.ExpectedOutput)
				result.Passed = outputTrimmed == expectedTrimmed
				if !result.Passed {
					result.Error = fmt.Sprintf("Expected '%s' but got '%s'", expectedTrimmed, outputTrimmed)
				}
			}
		case <-time.After(timeout):
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
			result.Error = "Time limit exceeded"
		}

		results = append(results, result)
	}

	return results
}

func allTestsPassed(results []TestResult) bool {
	for _, result := range results {
		if !result.Passed {
			return false
		}
	}
	return true
}
