package runner

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strathlearn/backend/db"
	"strconv"
	"strings"
	"time"

	"strathlearn/backend/models"
	"strathlearn/backend/utils"
)

const (
	judge0BaseURL = "http://172.17.0.1:2358"
	submitURL     = judge0BaseURL + "/submissions"
	statusURL     = judge0BaseURL + "/submissions/%s?base64_encoded=true"
)

type Judge0Runner struct {
	client *http.Client
}

type Judge0SubmissionRequest struct {
	SourceCode     string  `json:"source_code"`
	Language       int     `json:"language_id"`
	Stdin          string  `json:"stdin,omitempty"`
	ExpectedOutput string  `json:"expected_output,omitempty"`
	Base64Encoded  bool    `json:"base64_encoded"`
	CPUTimeLimit   float64 `json:"cpu_time_limit,omitempty"`
	MemoryLimit    int     `json:"memory_limit,omitempty"`
	CompileTimeout int     `json:"compile_timeout,omitempty"`
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
	return true
}

func (r *Judge0Runner) RunTests(code string, challenge models.Challenge) []models.TestResult {
	log.Printf("Running tests for challenge: %s", challenge.ID)
	results := make([]models.TestResult, 0, len(challenge.TestCases))

	for _, tc := range challenge.TestCases {
		result := models.TestResult{
			TestCaseID: tc.ID,
			Passed:     false,
			Output:     "",
			Error:      "",
		}

		token, err := r.submitCode(code, tc.Input, challenge.TimeLimit, challenge.MemoryLimit, challenge.ID, "")
		if err != nil {
			log.Printf("Error submitting code: %v", err)
			result.Error = fmt.Sprintf("Submission error: %v", err)
			results = append(results, result)
			continue
		}

		response, err := r.waitForResult(token)
		if err != nil {
			log.Printf("Error getting submission result: %v", err)
			result.Error = fmt.Sprintf("Execution error: %v", err)
			results = append(results, result)
			continue
		}

		switch response.Status.ID {
		case 3:
			result.Output = utils.CleanOutput(response.Stdout)

			if response.Time != "" {
				executionTime, err := strconv.ParseFloat(response.Time, 64)
				if err == nil {
					result.ExecutionTime = executionTime
				}
			}

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
		case 5:
			result.Error = "Time limit exceeded"
		case 6:
			result.Error = "Compilation error: " + response.CompileOutput
		case 11:
			result.Output = response.Stdout
			result.Error = "Runtime error: " + response.Message
		default:
			result.Error = fmt.Sprintf("Error: %s", response.Status.Description)
			if response.CompileOutput != "" {
				result.Error += " - " + response.CompileOutput
			}
			if response.Stderr != "" {
				result.Error += " - " + response.Stderr
			}
			if response.Message != "" {
				result.Error += " - " + response.Message
			}
		}

		results = append(results, result)
	}

	return results
}

func decodeBase64(input string) string {
	if input == "" {
		return ""
	}

	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		log.Printf("Error decoding base64: %v", err)
		return ""
	}

	return string(decoded)
}

func (r *Judge0Runner) submitCode(code, input string, timeLimit, memoryLimit int, challengeID, userID string) (string, error) {
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
	if err != nil {
		return "", fmt.Errorf("error marshaling submission: %w", err)
	}

	log.Printf("Submitting code to Judge0: %s", string(jsonData))

	req, err := http.NewRequest("POST", submitURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to submit code, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if response.Token == "" {
		return "", fmt.Errorf("received empty token from Judge0")
	}

	log.Printf("Successfully submitted code, token: %s", response.Token)
	return response.Token, nil
}

func (r *Judge0Runner) waitForResult(token string) (*Judge0SubmissionResponse, error) {
	url := fmt.Sprintf(statusURL, token)
	maxRetries := 20
	retryDelay := 500 * time.Millisecond

	log.Printf("Waiting for result from token: %s", token)

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating status request: %w", err)
		}

		resp, err := r.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error checking submission status: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error status response: %d, body: %s", resp.StatusCode, string(body))
		}

		var response Judge0SubmissionResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("error parsing response JSON: %w", err)
		}

		response.Stdout = decodeBase64(response.Stdout)
		response.Stderr = decodeBase64(response.Stderr)
		response.CompileOutput = decodeBase64(response.CompileOutput)
		response.Message = decodeBase64(response.Message)

		response.Stdout = strings.TrimSpace(response.Stdout)
		response.Stderr = strings.TrimSpace(response.Stderr)
		response.CompileOutput = strings.TrimSpace(response.CompileOutput)
		response.Message = strings.TrimSpace(response.Message)

		timeValue, _ := strconv.ParseFloat(response.Time, 64)
		dbSubmission := db.Submission{
			ID:            token,
			ChallengeID:   "",
			UserID:        "",
			Language:      "C",
			Code:          "",
			Stdout:        response.Stdout,
			Stderr:        response.Stderr,
			CompileOutput: response.CompileOutput,
			Message:       response.Message,
			StatusCode:    response.Status.ID,
			StatusDesc:    response.Status.Description,
			Memory:        response.Memory,
			Time:          timeValue,
			Token:         token,
			CreatedAt:     time.Now(),
		}

		if err := db.DB.Create(&dbSubmission).Error; err != nil {
			log.Printf("Failed to save submission to database: %v", err)
		}

		if response.Status.ID >= 3 {
			log.Printf("Submission completed with status: %d - %s",
				response.Status.ID, response.Status.Description)
			return &response, nil
		}

		log.Printf("Submission status: %d - %s, retrying in %v...",
			response.Status.ID, response.Status.Description, retryDelay)
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("timed out waiting for submission result after %d attempts", maxRetries)
}
