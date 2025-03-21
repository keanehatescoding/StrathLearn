package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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
	FilePath    string     `json:"-"`
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
var dockerClient *client.Client
var containerCleanupChannel = make(chan string, 100)

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
		go containerCleanupWorker()
	}

	challengesDir := "./backend/challenges"
	if _, err := os.Stat(challengesDir); os.IsNotExist(err) {
		log.Printf("Creating challenges directory at %s", challengesDir)
		tempDir, err := os.MkdirTemp(".", "challenges")
		if err != nil {
			log.Fatalf("Failed to create challenges directory: %v", err)
		}
		challengesDir = tempDir
		createSampleChallenge(challengesDir)
	}

	challenges = loadChallenges(challengesDir)
	log.Printf("Loaded %d challenges", len(challenges))

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

		fmt.Fprintf(w, "\nDocker connection status: %v\n", dockerClient != nil)
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

func containerCleanupWorker() {
	ctx := context.Background()
	for containerId := range containerCleanupChannel {
		time.Sleep(500 * time.Millisecond)
		err := dockerClient.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{Force: true})
		if err != nil && !strings.Contains(err.Error(), "No such container") {
			log.Printf("Error removing container %s: %v", containerId, err)
		}
	}
}

func scheduleContainerCleanup(containerId string) {
	containerCleanupChannel <- containerId
}

func createSampleChallenge(dir string) {
	challengeJSON := `{
        "id": "hello-world",
        "title": "Hello, World",
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
                "expectedOutput": "Hello, World!",
                "hidden": false
            }
        ],
        "initialCode": "#include <stdio.h>\\n\\nint main() {\\n    // Write your code here\\n    \\n    return 0;\\n}",
        "solutions": [
            "#include <stdio.h>\\n\\nint main() {\\n    printf(\\\"Hello, World!\\\");\\n    return 0;\\n}"
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

			var challenge Challenge
			if err := json.Unmarshal(data, &challenge); err != nil {
				log.Printf("JSON parse error in %s: %v", filePath, err)
				continue
			}

			challenge.FilePath = filePath

			if challenge.ID == "" {
				challenge.ID = strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			}

			log.Printf("Loaded challenge from %s: ID=%s, Title=%s",
				filePath, challenge.ID, challenge.Title)

			challenges[challenge.ID] = challenge
		}
	}

	if len(challenges) == 0 {
		log.Println("No challenges found, creating sample challenge")
		createSampleChallenge(dir)
		return loadChallenges(dir)
	}

	return challenges
}

func runTests(code string, challenge Challenge) []TestResult {
	log.Printf("Running tests for challenge: %s", challenge.ID)

	if dockerClient != nil {
		return runTestsInDocker(code, challenge)
	}

	log.Println("Using local execution mode")
	return runTestsLocally(code, challenge)
}

func runTestsInDocker(code string, challenge Challenge) []TestResult {
	results := make([]TestResult, 0, len(challenge.TestCases))
	submissionID := uuid.New().String()

	codeDir := "/code/" + submissionID
	sourcePath := filepath.Join(codeDir, "solution.c")

	code = strings.ReplaceAll(code, "\r\n", "\n")
	code = strings.ReplaceAll(code, "\r", "\n")

	ctx := context.Background()

	setupConfig := &container.Config{
		Image: "strathlearn-code-runner:latest",
		Cmd: []string{
			"sh", "-c",
			fmt.Sprintf("mkdir -p %s && cat > %s", codeDir, sourcePath),
		},
		Tty:          false,
		OpenStdin:    true,
		StdinOnce:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}

	hostConfig := &container.HostConfig{
		Binds: []string{"code-runner-data:/code"},
	}

	setupResp, err := dockerClient.ContainerCreate(
		ctx, setupConfig, hostConfig, nil, nil, "setup-"+submissionID)
	if err != nil {
		log.Printf("Error creating setup container: %v", err)
		return []TestResult{{
			TestCaseID: "setup",
			Passed:     false,
			Error:      "System error: could not create setup container",
		}}
	}

	attachResp, err := dockerClient.ContainerAttach(
		ctx, setupResp.ID, types.ContainerAttachOptions{
			Stdin:  true,
			Stdout: true,
			Stderr: true,
			Stream: true,
		})
	if err != nil {
		log.Printf("Error attaching to setup container: %v", err)
		scheduleContainerCleanup(setupResp.ID)
		return []TestResult{{
			TestCaseID: "setup",
			Passed:     false,
			Error:      "System error: could not attach to setup container",
		}}
	}
	defer attachResp.Close()

	if err := dockerClient.ContainerStart(
		ctx, setupResp.ID, types.ContainerStartOptions{}); err != nil {
		log.Printf("Error starting setup container: %v", err)
		scheduleContainerCleanup(setupResp.ID)
		return []TestResult{{
			TestCaseID: "setup",
			Passed:     false,
			Error:      "System error: could not start setup container",
		}}
	}

	go func() {
		io.Copy(attachResp.Conn, strings.NewReader(code))
		attachResp.CloseWrite()
	}()

	statusCh, errCh := dockerClient.ContainerWait(
		ctx, setupResp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Printf("Error waiting for setup container: %v", err)
			scheduleContainerCleanup(setupResp.ID)
			return []TestResult{{
				TestCaseID: "setup",
				Passed:     false,
				Error:      "System error: setup container error",
			}}
		}
	case <-statusCh:
	}

	logs, err := dockerClient.ContainerLogs(
		ctx, setupResp.ID, types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		})
	if err == nil {
		defer logs.Close()
		var logBuf bytes.Buffer
		io.Copy(&logBuf, logs)
		log.Printf("Setup container logs: %s", logBuf.String())
	}

	scheduleContainerCleanup(setupResp.ID)

	compileConfig := &container.Config{
		Image: "strathlearn-code-runner:latest",
		Cmd: []string{
			"sh", "-c",
			fmt.Sprintf("cd %s && gcc -Wall -o solution solution.c", codeDir),
		},
		Tty: false,
	}

	compileResp, err := dockerClient.ContainerCreate(
		ctx, compileConfig, hostConfig, nil, nil, "compile-"+submissionID)
	if err != nil {
		log.Printf("Error creating compile container: %v", err)
		return []TestResult{{
			TestCaseID: "compile",
			Passed:     false,
			Error:      "System error: could not create compile container",
		}}
	}

	if err := dockerClient.ContainerStart(
		ctx, compileResp.ID, types.ContainerStartOptions{}); err != nil {
		log.Printf("Error starting compile container: %v", err)
		scheduleContainerCleanup(compileResp.ID)
		return []TestResult{{
			TestCaseID: "compile",
			Passed:     false,
			Error:      "System error: could not start compile container",
		}}
	}

	statusCh, errCh = dockerClient.ContainerWait(
		ctx, compileResp.ID, container.WaitConditionNotRunning)
	var compileExitCode int64
	select {
	case err := <-errCh:
		if err != nil {
			log.Printf("Error waiting for compile container: %v", err)
			scheduleContainerCleanup(compileResp.ID)
			return []TestResult{{
				TestCaseID: "compile",
				Passed:     false,
				Error:      "System error: compile container error",
			}}
		}
	case status := <-statusCh:
		compileExitCode = status.StatusCode
	}

	compileOutput, err := dockerClient.ContainerLogs(
		ctx, compileResp.ID, types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		})

	var compileErrorMsg string
	if err == nil {
		defer compileOutput.Close()
		var outputBuf bytes.Buffer
		_, err = stdCopy(&outputBuf, &outputBuf, compileOutput)
		if err != nil {
			log.Printf("Error reading compile logs: %v", err)
		} else if compileExitCode != 0 {
			compileErrorMsg = outputBuf.String()
			scheduleContainerCleanup(compileResp.ID)
			return []TestResult{{
				TestCaseID: "compile",
				Passed:     false,
				Error:      "Compilation error: " + compileErrorMsg,
			}}
		}
	}

	scheduleContainerCleanup(compileResp.ID)

	for _, tc := range challenge.TestCases {
		result := TestResult{
			TestCaseID: tc.ID,
			Passed:     false,
		}

		inputPath := filepath.Join(codeDir, "input.txt")
		inputConfig := &container.Config{
			Image: "strathlearn-code-runner:latest",
			Cmd: []string{
				"sh", "-c",
				fmt.Sprintf("echo -n %s > %s",
					shellescape(tc.Input),
					inputPath),
			},
			Tty: false,
		}

		inputResp, err := dockerClient.ContainerCreate(
			ctx, inputConfig, hostConfig, nil, nil, "input-"+submissionID+"-"+tc.ID)
		if err != nil {
			log.Printf("Error creating input container: %v", err)
			result.Error = "System error: failed to prepare input"
			results = append(results, result)
			continue
		}

		if err := dockerClient.ContainerStart(
			ctx, inputResp.ID, types.ContainerStartOptions{}); err != nil {
			log.Printf("Error starting input container: %v", err)
			scheduleContainerCleanup(inputResp.ID)
			result.Error = "System error: failed to prepare input"
			results = append(results, result)
			continue
		}

		dockerClient.ContainerWait(
			ctx, inputResp.ID, container.WaitConditionNotRunning)
		scheduleContainerCleanup(inputResp.ID)

		runConfig := &container.Config{
			Image: "strathlearn-code-runner:latest",
			Cmd: []string{
				"sh", "-c",
				fmt.Sprintf("cd %s && timeout %d ./solution < input.txt",
					codeDir,
					challenge.TimeLimit),
			},
			Tty: false,
		}

		runHostConfig := &container.HostConfig{
			Resources: container.Resources{
				Memory:     int64(challenge.MemoryLimit) * 1024 * 1024,
				MemorySwap: int64(challenge.MemoryLimit) * 1024 * 1024,
				CPUPeriod:  100000,
				CPUQuota:   50000,
			},
			Binds: []string{"code-runner-data:/code"},
		}

		runResp, err := dockerClient.ContainerCreate(
			ctx, runConfig, runHostConfig, nil, nil, "run-"+submissionID+"-"+tc.ID)
		if err != nil {
			log.Printf("Error creating run container: %v", err)
			result.Error = "System error: could not create run container"
			results = append(results, result)
			continue
		}

		if err := dockerClient.ContainerStart(
			ctx, runResp.ID, types.ContainerStartOptions{}); err != nil {
			log.Printf("Error starting run container: %v", err)
			scheduleContainerCleanup(runResp.ID)
			result.Error = "System error: could not start run container"
			results = append(results, result)
			continue
		}

		timeoutCtx, cancel := context.WithTimeout(
			ctx, time.Duration(challenge.TimeLimit+1)*time.Second)
		defer cancel()

		var runExitCode int64
		select {
		case <-timeoutCtx.Done():
			dockerClient.ContainerStop(ctx, runResp.ID, container.StopOptions{})
			result.Error = "Time limit exceeded"
		case err := <-func() <-chan error {
			ch := make(chan error, 1)
			go func() {
				statusCh, errCh := dockerClient.ContainerWait(
					ctx, runResp.ID, container.WaitConditionNotRunning)
				select {
				case err := <-errCh:
					ch <- err
				case status := <-statusCh:
					runExitCode = status.StatusCode
					ch <- nil
				}
			}()
			return ch
		}():
			if err != nil {
				log.Printf("Error waiting for run container: %v", err)
				result.Error = "System error: error waiting for program execution"
			} else if runExitCode == 124 {
				result.Error = "Time limit exceeded"
			} else if runExitCode != 0 {
				result.Error = fmt.Sprintf("Runtime error: program exited with code %d", runExitCode)
			}
		}

		output, err := dockerClient.ContainerLogs(
			ctx, runResp.ID, types.ContainerLogsOptions{
				ShowStdout: true,
				ShowStderr: true,
			})
		if err == nil {
			defer output.Close()
			var outputBuf bytes.Buffer
			_, err = stdCopy(&outputBuf, &outputBuf, output)
			if err != nil {
				log.Printf("Error reading output: %v", err)
				result.Error = "System error: failed to read program output"
			} else {
				programOutput := cleanOutput(outputBuf.String())
				result.Output = programOutput

				if result.Error == "" {
					expectedOutput := cleanOutput(tc.ExpectedOutput)

					log.Printf("Expected output: '%s'", expectedOutput)
					log.Printf("Actual output: '%s'", programOutput)

					result.Passed = programOutput == expectedOutput
					if !result.Passed {
						result.Error = fmt.Sprintf("Expected '%s' but got '%s'",
							formatForDisplay(expectedOutput),
							formatForDisplay(programOutput))
					}
				}
			}
		} else {
			log.Printf("Error getting run logs: %v", err)
			result.Error = "System error: failed to get program output"
		}

		scheduleContainerCleanup(runResp.ID)
		results = append(results, result)
	}

	return results
}

func runTestsLocally(code string, challenge Challenge) []TestResult {
	results := make([]TestResult, 0, len(challenge.TestCases))

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

		timeout := time.Duration(challenge.TimeLimit) * time.Second
		if timeout == 0 {
			timeout = 5 * time.Second
		}

		outputChan := make(chan []byte, 1)
		errorChan := make(chan error, 1)

		go func() {
			output, err := cmd.CombinedOutput()
			outputChan <- output
			errorChan <- err
		}()

		select {
		case output := <-outputChan:
			err := <-errorChan
			if err != nil {
				result.Error = "Runtime error: " + err.Error()
			} else {
				programOutput := cleanOutput(string(output))
				result.Output = programOutput

				expectedOutput := cleanOutput(tc.ExpectedOutput)

				result.Passed = programOutput == expectedOutput
				if !result.Passed {
					result.Error = fmt.Sprintf("Expected '%s' but got '%s'",
						formatForDisplay(expectedOutput),
						formatForDisplay(programOutput))
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
	return len(results) > 0
}

func shellescape(s string) string {
	if s == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func cleanOutput(output string) string {
	output = strings.ReplaceAll(output, "\r\n", "\n")
	output = strings.ReplaceAll(output, "\r", "\n")

	startIndex := 0
	for i, c := range output {
		if c >= 32 && c <= 126 {
			startIndex = i
			break
		}
	}
	if startIndex > 0 {
		output = output[startIndex:]
	}

	return strings.TrimSpace(output)
}

func formatForDisplay(s string) string {
	return strings.ReplaceAll(s, "\n", "\\n")
}

func stdCopy(stdout, stderr io.Writer, src io.Reader) (written int64, err error) {
	var buf [8]byte
	var totalN int64

	for {
		_, err = io.ReadFull(src, buf[:])
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return totalN, err
		}

		frameSize := int64(buf[4])<<24 | int64(buf[5])<<16 | int64(buf[6])<<8 | int64(buf[7])

		var dst io.Writer
		if buf[0] == 1 {
			dst = stdout
		} else {
			dst = stderr
		}

		n, err := io.CopyN(dst, src, frameSize)
		totalN += n
		if err != nil {
			return totalN, err
		}
	}
}
