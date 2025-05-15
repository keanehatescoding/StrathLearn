package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"strathlearn/backend/models"
	"strathlearn/backend/utils"
)

const (
	judge0BaseURL = "http://host.docker.internal:2358"
	submitURL     = judge0BaseURL + "/submissions"
	statusURL     = judge0BaseURL + "/submissions/%s"
)

type Judge0Runner struct {
	client *http.Client
}

type Judge0SubmissionRequest struct {
	SourceCode     string  `json:"source_code"`
	Language       int     `json:"language_id"` // 103 for C (gcc)
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
		result := models.TestResult{
			TestCaseID: tc.ID,
			Passed:     false,
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
		Language:       50,
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

	return response.Token, nil
}

func (r *Judge0Runner) waitForResult(token string) (*Judge0SubmissionResponse, error) {
	url := fmt.Sprintf(statusURL, token)

	// Poll until we get a result
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

		// Check if processing is complete
		if response.Status.ID >= 3 {
			return &response, nil
		}

		// Wait before polling again
		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("timed out waiting for submission result")
}
