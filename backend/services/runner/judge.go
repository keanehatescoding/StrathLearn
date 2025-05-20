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

type Judge0Runner struct {
	client *http.Client
}

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

type Judge0Status struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

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

const (
	LangC          = 50
	LangCPP        = 54
	LangJava       = 62
	LangPython     = 71
	LangJavaScript = 63
)

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
	log.Printf("Running tests for challenge: %s with %d test cases", challenge.ID, len(challenge.TestCases))
	results := make([]models.TestResult, 0, len(challenge.TestCases))

	for _, tc := range challenge.TestCases {
		result := models.TestResult{
			TestCaseID:    tc.ID,
			Passed:        false,
			Output:        "",
			Error:         "",
			ExecutionTime: 0,
			Memory:        0,
		}

		token, err := r.submitCode(code, tc.Input, challenge.TimeLimit, challenge.MemoryLimit, challenge.ID)
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

		result = r.processResponse(response, tc.ExpectedOutput, result)

		results = append(results, result)
	}

	return results
}

func (r *Judge0Runner) processResponse(response *Judge0SubmissionResponse, expectedOutput string, result models.TestResult) models.TestResult {
	if response.Time != "" {
		executionTime, err := strconv.ParseFloat(response.Time, 64)
		if err == nil {
			result.ExecutionTime = executionTime
		}
	}

	result.Memory = response.Memory

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

func (r *Judge0Runner) submitCode(code, input string, timeLimit, memoryLimit int, challengeID string) (string, error) {
	if timeLimit <= 0 {
		timeLimit = 5
	}

	if memoryLimit <= 0 {
		memoryLimit = 512 * 1024
	} else {
		memoryLimit = memoryLimit * 1024
	}

	submission := Judge0SubmissionRequest{
		SourceCode:           code,
		LanguageID:           LangC,
		Stdin:                input,
		CPUTimeLimit:         float64(timeLimit),
		CPUExtraTime:         1.0,
		WallTimeLimit:        float64(timeLimit) + 2.0,
		MemoryLimit:          memoryLimit,
		StackLimit:           memoryLimit,
		MaxProcessesThreads:  30,
		EnablePerThreadLimit: true,
		EnablePerThreadMem:   true,
		MaxFileSize:          1024,
		RedirectStderr:       false,
		EnableNetwork:        false,
	}

	jsonData, err := json.Marshal(submission)
	if err != nil {
		return "", fmt.Errorf("error marshaling submission request: %v", err)
	}

	log.Printf("Submitting code to Judge0: %s", string(jsonData))

	req, err := http.NewRequest("POST", submitURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("base64_encoded", "false")
	req.URL.RawQuery = q.Encode()

	resp, err := r.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to submit code, status: %d, response: %s",
			resp.StatusCode, string(bodyBytes))
	}

	var response Judge0SubmissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	timeValue, _ := strconv.ParseFloat(response.Time, 64)

	dbSubmission := db.Submission{
		ID:            response.Token,
		ChallengeID:   challengeID,
		UserID:        "",
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
	}

	return response.Token, nil
}

func (r *Judge0Runner) waitForResult(token string) (*Judge0SubmissionResponse, error) {
	url := fmt.Sprintf(statusURL, token)

	maxRetries := 60
	retryInterval := 500 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating status request: %v", err)
		}

		q := req.URL.Query()
		q.Add("base64_encoded", "false")
		q.Add("fields", "*")
		req.URL.RawQuery = q.Encode()

		resp, err := r.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error checking submission status: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("error response from Judge0: %d", resp.StatusCode)
		}

		var response Judge0SubmissionResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("error decoding status response: %v", err)
		}

		if response.Status.ID >= 3 {
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
			}

			return &response, nil
		}

		time.Sleep(retryInterval)

		if i > 10 {
			retryInterval = 1 * time.Second
		}
	}

	return nil, fmt.Errorf("timed out waiting for submission result after %d attempts", maxRetries)
}

func (r *Judge0Runner) GetSubmission(token string) (*Judge0SubmissionResponse, error) {
	url := fmt.Sprintf(statusURL, token)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	q := req.URL.Query()
	q.Add("base64_encoded", "false")
	q.Add("fields", "*")
	req.URL.RawQuery = q.Encode()

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching submission: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from Judge0: %d", resp.StatusCode)
	}

	var response Judge0SubmissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &response, nil
}

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
