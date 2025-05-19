package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strathlearn/backend/db"
	"strconv"
	"time"

	"strathlearn/backend/models"
	"strathlearn/backend/utils"
)

const (
	judge0BaseURL = "http://172.17.0.1:2358"
	submitURL     = judge0BaseURL + "/submissions"
	statusURL     = judge0BaseURL + "/submissions/%s"
)

type Judge0Runner struct {
	client *http.Client
}

type Judge0SubmissionRequest struct {
	SourceCode     string  `json:"source_code"`
	Language       int     `json:"language_id"` // Language ID (50 for C, as seen in frontend)
	Stdin          string  `json:"stdin,omitempty"`
	ExpectedOutput string  `json:"expected_output,omitempty"`
	Base64Encoded  bool    `json:"base64_encoded"`
	CPUTimeLimit   float64 `json:"cpu_time_limit,omitempty"`  // seconds
	MemoryLimit    int     `json:"memory_limit,omitempty"`    // KB
	CompileTimeout int     `json:"compile_timeout,omitempty"` // seconds
}

type Judge0Status struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type Judge0SubmissionResponse struct {
	Stdout        string       `json:"stdout"`
	Time          string       `json:"time"`
	Memory        int          `json:"memory"`
	Stderr        string       `json:"stderr"`
	Token         string       `json:"token"`
	CompileOutput string       `json:"compile_output"`
	Message       string       `json:"message"`
	Status        Judge0Status `json:"status"`
}

func NewJudge0Runner() *Judge0Runner {
	return &Judge0Runner{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (r *Judge0Runner) IsDockerAvailable() bool {
	// Keep the same interface as the old runner
	return true
}

func (r *Judge0Runner) RunTests(code string, challenge models.Challenge) []models.TestResult {
	log.Printf("Running tests for challenge: %s", challenge.ID)
	results := make([]models.TestResult, 0, len(challenge.TestCases))

	for _, tc := range challenge.TestCases {
		// Initialize result with expected frontend keys
		result := models.TestResult{
			TestCaseID: tc.ID, // Frontend expects camelCase
			Passed:     false,
			Output:     "", // Initialize empty string to avoid null in JSON
			Error:      "", // Initialize empty string to avoid null in JSON
		}

		token, err := r.submitCode(code, tc.Input, challenge.TimeLimit, challenge.MemoryLimit)
		if err != nil {
			log.Printf("Error submitting code: %v", err)
			result.Error = fmt.Sprintf("Submission error: %v", err)
			results = append(results, result)
			continue
		}

		// Poll for submission status
		response, err := r.waitForResult(token)
		if err != nil {
			log.Printf("Error getting submission result: %v", err)
			result.Error = fmt.Sprintf("Execution error: %v", err)
			results = append(results, result)
			continue
		}

		// Process the response
		if response.Status.ID == 3 { // 3 means "Accepted" in Judge0
			result.Output = utils.CleanOutput(response.Stdout)

			// Set execution time as float
			if response.Time != "" {
				executionTime, err := strconv.ParseFloat(response.Time, 64)
				if err == nil {
					result.ExecutionTime = executionTime
				}
			}

			// Set memory usage
			result.Memory = response.Memory

			expectedOutput := utils.CleanOutput(tc.ExpectedOutput)
			log.Printf("Expected output: '%s'", expectedOutput)
			log.Printf("Actual output: '%s'", result.Output)

			result.Passed = result.Output == expectedOutput
			if !result.Passed {
				result.Error = fmt.Sprintf("Expected '%s' but got '%s'",
					utils.FormatForDisplay(expectedOutput),
					utils.FormatForDisplay(result.Output))
			}
		} else if response.CompileOutput != "" {
			result.Error = "Compilation error: " + response.CompileOutput
		} else if response.Status.ID == 5 { // Time Limit Exceeded
			result.Error = "Time limit exceeded"
		} else {
			result.Error = fmt.Sprintf("Runtime error: %s", response.Status.Description)
			if response.Stderr != "" {
				result.Error += " - " + response.Stderr
			}
		}

		results = append(results, result)
	}

	return results
}

func (r *Judge0Runner) submitCode(code, input string, timeLimit, memoryLimit int) (string, error) {
	submission := Judge0SubmissionRequest{
		SourceCode:     code,
		Language:       50, // Match the frontend's default language ID for C
		Stdin:          input,
		Base64Encoded:  false,
		CPUTimeLimit:   float64(timeLimit),
		MemoryLimit:    memoryLimit * 1024,
		CompileTimeout: 10,
	}

	jsonData, err := json.Marshal(submission)
	fmt.Print("Submitting code: ", string(jsonData))
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", submitURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to submit code, status: %d", resp.StatusCode)
	}

	var response Judge0SubmissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	timeValue, _ := strconv.ParseFloat(response.Time, 64)
	dbSubmission := db.Submission{
		ID:            response.Token, // Using the token as the ID
		ChallengeID:   "",             // You'll need to pass the challenge ID to this function
		UserID:        "",             // You'll need to pass the user ID to this function
		Language:      "C",            // Based on language ID 50, you might want to store actual language name
		Code:          code,
		Stdout:        response.Stdout,
		Stderr:        response.Stderr,
		CompileOutput: response.CompileOutput,
		Message:       response.Message,
		StatusCode:    response.Status.ID,
		StatusDesc:    response.Status.Description,
		Memory:        response.Memory,
		Time:          timeValue,
		Token:         response.Token,
	}

	if err := db.DB.Create(&dbSubmission).Error; err != nil {
		log.Printf("Failed to save submission to database: %v", err)
	}

	return response.Token, nil
}

func (r *Judge0Runner) waitForResult(token string) (*Judge0SubmissionResponse, error) {
	url := fmt.Sprintf(statusURL, token)

	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := r.client.Do(req)
		if err != nil {
			return nil, err
		}

		var response Judge0SubmissionResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		if response.Status.ID >= 3 {
			return &response, nil
		}

		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("timed out waiting for submission result")
}
