package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

// Judge0Runner implements the Runner interface using Judge0 API
type Judge0Runner struct {
	client *http.Client
}

// Judge0SubmissionRequest represents the request body for creating a submission
type Judge0SubmissionRequest struct {
	SourceCode           string  `json:"source_code"`
	LanguageID           int     `json:"language_id"`
	Stdin                string  `json:"stdin,omitempty"`
	ExpectedOutput       string  `json:"expected_output,omitempty"`
	CPUTimeLimit         float64 `json:"cpu_time_limit,omitempty"`
	CPUExtraTime         float64 `json:"cpu_extra_time,omitempty"`
	WallTimeLimit        float64 `json:"wall_time_limit,omitempty"`
	MemoryLimit          int     `json:"memory_limit,omitempty"`
	StackLimit           int     `json:"stack_limit,omitempty"`
	MaxProcessesThreads  int     `json:"max_processes_and_or_threads,omitempty"`
	EnablePerThreadLimit bool    `json:"enable_per_process_and_thread_time_limit,omitempty"`
	EnablePerThreadMem   bool    `json:"enable_per_process_and_thread_memory_limit,omitempty"`
	MaxFileSize          int     `json:"max_file_size,omitempty"`
	RedirectStderr       bool    `json:"redirect_stderr_to_stdout,omitempty"`
	EnableNetwork        bool    `json:"enable_network,omitempty"`
	CompilerOptions      string  `json:"compiler_options,omitempty"`
	CommandLineArgs      string  `json:"command_line_arguments,omitempty"`
}

// Judge0Status represents the status of a submission
type Judge0Status struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

// Judge0SubmissionResponse represents the response from Judge0 API
type Judge0SubmissionResponse struct {
	Stdout        string       `json:"stdout"`
	Stderr        string       `json:"stderr"`
	CompileOutput string       `json:"compile_output"`
	Message       string       `json:"message"`
	ExitCode      int          `json:"exit_code"`
	ExitSignal    int          `json:"exit_signal"`
	Status        Judge0Status `json:"status"`
	CreatedAt     string       `json:"created_at"`
	FinishedAt    string       `json:"finished_at"`
	Token         string       `json:"token"`
	Time          string       `json:"time"`
	WallTime      string       `json:"wall_time"`
	Memory        int          `json:"memory"`
}

// Status ID constants from Judge0
const (
	StatusInQueue             = 1
	StatusProcessing          = 2
	StatusAccepted            = 3
	StatusWrongAnswer         = 4
	StatusTimeLimitExceeded   = 5
	StatusCompilationError    = 6
	StatusRuntimeErrorSIGSEGV = 7
	StatusRuntimeErrorSIGXFSZ = 8
	StatusRuntimeErrorSIGFPE  = 9
	StatusRuntimeErrorSIGABRT = 10
	StatusRuntimeErrorNZEC    = 11
	StatusRuntimeErrorOther   = 12
	StatusInternalError       = 13
	StatusExecFormatError     = 14
)

// Language ID constants for Judge0
const (
	LangC          = 50
	LangCPP        = 54
	LangJava       = 62
	LangPython     = 71
	LangJavaScript = 63
)

// NewJudge0Runner creates a new Judge0Runner
func NewJudge0Runner() *Judge0Runner {
	return &Judge0Runner{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsDockerAvailable implements the Runner interface
func (r *Judge0Runner) IsDockerAvailable() bool {
	// Always return true as Judge0 is running in a separate container
	return true
}

// RunTests runs the tests for a challenge
func (r *Judge0Runner) RunTests(code string, challenge models.Challenge) []models.TestResult {
	log.Printf("Running tests for challenge: %s with %d test cases", challenge.ID, len(challenge.TestCases))
	results := make([]models.TestResult, 0, len(challenge.TestCases))

	for _, tc := range challenge.TestCases {
		// Initialize result with expected frontend keys
		result := models.TestResult{
			TestCaseID:    tc.ID,
			Passed:        false,
			Output:        "",
			Error:         "",
			ExecutionTime: 0,
			Memory:        0,
		}

		// Submit code to Judge0
		token, err := r.submitCode(code, tc.Input, challenge.TimeLimit, challenge.MemoryLimit, challenge.ID)
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

		// Process the response based on status
		result = r.processResponse(response, tc.ExpectedOutput, result)

		results = append(results, result)
	}

	return results
}

// processResponse processes the Judge0 response and updates the test result
func (r *Judge0Runner) processResponse(response *Judge0SubmissionResponse, expectedOutput string, result models.TestResult) models.TestResult {
	// Set execution time if available
	if response.Time != "" {
		executionTime, err := strconv.ParseFloat(response.Time, 64)
		if err == nil {
			result.ExecutionTime = executionTime
		}
	}

	// Set memory usage
	result.Memory = response.Memory

	// Process based on status
	switch response.Status.ID {
	case StatusAccepted:
		result.Output = utils.CleanOutput(response.Stdout)
		expectedOutput = utils.CleanOutput(expectedOutput)

		log.Printf("Expected output: '%s'", expectedOutput)
		log.Printf("Actual output: '%s'", result.Output)

		result.Passed = result.Output == expectedOutput
		if !result.Passed {
			result.Error = fmt.Sprintf("Expected '%s' but got '%s'",
				utils.FormatForDisplay(expectedOutput),
				utils.FormatForDisplay(result.Output))
		}

	case StatusCompilationError:
		result.Error = "Compilation error: " + response.CompileOutput

	case StatusTimeLimitExceeded:
		result.Error = "Time limit exceeded"

	case StatusRuntimeErrorSIGSEGV:
		result.Error = "Runtime error: Segmentation fault"
		result.Output = response.Stdout

	case StatusRuntimeErrorSIGXFSZ:
		result.Error = "Runtime error: File size limit exceeded"
		result.Output = response.Stdout

	case StatusRuntimeErrorSIGFPE:
		result.Error = "Runtime error: Floating point error"
		result.Output = response.Stdout

	case StatusRuntimeErrorSIGABRT:
		result.Error = "Runtime error: Aborted"
		result.Output = response.Stdout

	case StatusRuntimeErrorNZEC:
		result.Error = "Runtime error: Non-zero exit code " + strconv.Itoa(response.ExitCode)
		result.Output = response.Stdout

	case StatusRuntimeErrorOther:
		result.Error = "Runtime error: " + response.Message
		result.Output = response.Stdout

	case StatusInternalError:
		result.Error = "Internal server error: " + response.Message

	case StatusExecFormatError:
		result.Error = "Execution format error"

	default:
		result.Error = fmt.Sprintf("Unknown error (Status %d): %s",
			response.Status.ID, response.Status.Description)
		if response.Message != "" {
			result.Error += " - " + response.Message
		}
		if response.Stderr != "" {
			result.Error += " - " + response.Stderr
		}
		result.Output = response.Stdout
	}

	return result
}

// submitCode submits code to Judge0
func (r *Judge0Runner) submitCode(code, input string, timeLimit, memoryLimit int, challengeID string) (string, error) {
	// Default to 2 seconds if timeLimit is 0
	if timeLimit <= 0 {
		timeLimit = 2
	}

	// Default to 128MB if memoryLimit is 0
	if memoryLimit <= 0 {
		memoryLimit = 128 * 1024 // 128MB in KB
	} else {
		memoryLimit = memoryLimit * 1024 // Convert MB to KB
	}

	// Keep memory limit within reasonable bounds (max 512MB)
	if memoryLimit > 512*1024 {
		memoryLimit = 512 * 1024
	}

	// Create submission request with minimal required parameters
	submission := Judge0SubmissionRequest{
		SourceCode:   code,
		LanguageID:   LangC, // Default to C
		Stdin:        input,
		CPUTimeLimit: float64(timeLimit),
		MemoryLimit:  memoryLimit,
		// Don't set stack_limit - use Judge0 default
	}

	jsonData, err := json.Marshal(submission)
	if err != nil {
		return "", fmt.Errorf("error marshaling submission request: %v", err)
	}

	log.Printf("Submitting code to Judge0: %s", string(jsonData))

	// Create the request
	req, err := http.NewRequest("POST", submitURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add query parameter for base64 encoding
	q := req.URL.Query()
	q.Add("base64_encoded", "false")
	req.URL.RawQuery = q.Encode()

	// Send the request
	resp, err := r.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to submit code, status: %d, response: %s",
			resp.StatusCode, string(bodyBytes))
	}

	// Parse the response
	var response Judge0SubmissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	// Store submission in database
	timeValue, _ := strconv.ParseFloat(response.Time, 64)

	dbSubmission := db.Submission{
		ID:            response.Token,
		ChallengeID:   challengeID,
		UserID:        "", // Would need to be passed in
		Language:      "C",
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
		CreatedAt:     time.Now(),
	}

	if err := db.DB.Create(&dbSubmission).Error; err != nil {
		log.Printf("Failed to save submission to database: %v", err)
		// Continue anyway, this is not critical
	}

	return response.Token, nil
}

// waitForResult waits for the submission result
func (r *Judge0Runner) waitForResult(token string) (*Judge0SubmissionResponse, error) {
	url := fmt.Sprintf(statusURL, token)

	maxRetries := 60                        // Increased for longer-running submissions
	retryInterval := 500 * time.Millisecond // 0.5 seconds

	for i := 0; i < maxRetries; i++ {
		// Create the request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating status request: %v", err)
		}

		// Add query parameter for base64 encoding
		q := req.URL.Query()
		q.Add("base64_encoded", "false")
		q.Add("fields", "*") // Get all fields
		req.URL.RawQuery = q.Encode()

		// Send the request
		resp, err := r.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error checking submission status: %v", err)
		}

		// Check response status
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("error response from Judge0: %d", resp.StatusCode)
		}

		// Parse the response
		var response Judge0SubmissionResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("error decoding status response: %v", err)
		}

		// Check if processing is complete (status >= 3)
		if response.Status.ID >= 3 {
			// Update the submission in the database with final status
			updateData := map[string]interface{}{
				"stdout":         response.Stdout,
				"stderr":         response.Stderr,
				"compile_output": response.CompileOutput,
				"message":        response.Message,
				"status_code":    response.Status.ID,
				"status_desc":    response.Status.Description,
				"memory":         response.Memory,
			}

			if response.Time != "" {
				timeValue, err := strconv.ParseFloat(response.Time, 64)
				if err == nil {
					updateData["time"] = timeValue
				}
			}

			if err := db.DB.Model(&db.Submission{}).Where("token = ?", token).Updates(updateData).Error; err != nil {
				log.Printf("Failed to update submission in database: %v", err)
				// Continue anyway, this is not critical
			}

			return &response, nil
		}

		// If still in queue or processing, wait and retry
		time.Sleep(retryInterval)

		// Increase retry interval after several attempts
		if i > 10 {
			retryInterval = 1 * time.Second
		}
	}

	return nil, fmt.Errorf("timed out waiting for submission result after %d attempts", maxRetries)
}

// GetSubmission gets a submission by token
func (r *Judge0Runner) GetSubmission(token string) (*Judge0SubmissionResponse, error) {
	url := fmt.Sprintf(statusURL, token)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add query parameter for base64 encoding
	q := req.URL.Query()
	q.Add("base64_encoded", "false")
	q.Add("fields", "*") // Get all fields
	req.URL.RawQuery = q.Encode()

	// Send the request
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching submission: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from Judge0: %d", resp.StatusCode)
	}

	// Parse the response
	var response Judge0SubmissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &response, nil
}

// GetJudge0Status returns a human-readable status message for a Judge0 status ID
func GetJudge0Status(statusID int) string {
	switch statusID {
	case StatusInQueue:
		return "In Queue"
	case StatusProcessing:
		return "Processing"
	case StatusAccepted:
		return "Accepted"
	case StatusWrongAnswer:
		return "Wrong Answer"
	case StatusTimeLimitExceeded:
		return "Time Limit Exceeded"
	case StatusCompilationError:
		return "Compilation Error"
	case StatusRuntimeErrorSIGSEGV:
		return "Runtime Error (SIGSEGV)"
	case StatusRuntimeErrorSIGXFSZ:
		return "Runtime Error (SIGXFSZ)"
	case StatusRuntimeErrorSIGFPE:
		return "Runtime Error (SIGFPE)"
	case StatusRuntimeErrorSIGABRT:
		return "Runtime Error (SIGABRT)"
	case StatusRuntimeErrorNZEC:
		return "Runtime Error (NZEC)"
	case StatusRuntimeErrorOther:
		return "Runtime Error (Other)"
	case StatusInternalError:
		return "Internal Error"
	case StatusExecFormatError:
		return "Execution Format Error"
	default:
		return fmt.Sprintf("Unknown Status (%d)", statusID)
	}
}
