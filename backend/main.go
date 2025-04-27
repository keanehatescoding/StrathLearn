package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
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

type ExecutionRequest struct {
	Code         string
	Challenge    Challenge
	SubmissionID string
}

type ExecutionResponse struct {
	Results []TestResult
	Error   string
}

type CodeExecutor struct {
	dockerClient *client.Client
	workDir      string
	mu           sync.Mutex
	containers   map[string]time.Time
}

func NewCodeExecutor(workDir string) (*CodeExecutor, error) {
	dockerClient, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
		client.WithHost("unix:///var/run/docker.sock"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Docker: %v", err)
	}

	if err := os.MkdirAll(workDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create work directory: %v", err)
	}

	executor := &CodeExecutor{
		dockerClient: dockerClient,
		workDir:      workDir,
		containers:   make(map[string]time.Time),
	}

	go executor.cleanupContainers()

	return executor, nil
}

func (e *CodeExecutor) Execute(req *ExecutionRequest, resp *ExecutionResponse) error {
	ctx := context.Background()
	submissionID := req.SubmissionID
	challenge := req.Challenge
	code := req.Code

	submissionDir := filepath.Join(e.workDir, submissionID)
	if err := os.MkdirAll(submissionDir, 0755); err != nil {
		resp.Error = fmt.Sprintf("Failed to create submission directory: %v", err)
		return nil
	}
	defer os.RemoveAll(submissionDir)

	sourcePath := filepath.Join(submissionDir, "solution.c")
	if err := os.WriteFile(sourcePath, []byte(code), 0644); err != nil {
		resp.Error = fmt.Sprintf("Failed to write source code: %v", err)
		return nil
	}

	seccompPath := filepath.Join(submissionDir, "seccomp.json")
	if err := os.WriteFile(seccompPath, []byte(seccompProfile), 0644); err != nil {
		resp.Error = fmt.Sprintf("Failed to write seccomp profile: %v", err)
		return nil
	}

	compileResult, err := e.compileCode(ctx, submissionID, submissionDir)
	if err != nil {
		resp.Error = fmt.Sprintf("System error during compilation: %v", err)
		return nil
	}
	if compileResult != "" {
		resp.Results = []TestResult{{
			TestCaseID: "compile",
			Passed:     false,
			Error:      "Compilation error: " + compileResult,
		}}
		return nil
	}

	results := make([]TestResult, 0, len(challenge.TestCases))
	for _, tc := range challenge.TestCases {
		result, err := e.runTest(ctx, submissionID, submissionDir, tc, challenge.TimeLimit, challenge.MemoryLimit)
		if err != nil {
			log.Printf("Error running test %s: %v", tc.ID, err)
			results = append(results, TestResult{
				TestCaseID: tc.ID,
				Passed:     false,
				Error:      fmt.Sprintf("System error: %v", err),
			})
			continue
		}
		results = append(results, result)
	}

	resp.Results = results
	return nil
}

func (e *CodeExecutor) compileCode(ctx context.Context, submissionID, submissionDir string) (string, error) {
	containerName := "compile-" + submissionID

	config := &container.Config{
		Image: "code-challenge-runner",
		Cmd: []string{
			"sh", "-c",
			"cd /code && gcc -Wall -Werror -std=c11 -o solution solution.c",
		},
		Tty:          false,
		AttachStdout: true,
		AttachStderr: true,
	}

	hostConfig := &container.HostConfig{
		AutoRemove: false,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: submissionDir,
				Target: "/code",
			},
		},
		Resources: container.Resources{
			Memory:     256 * 1024 * 1024, // 256MB
			MemorySwap: 256 * 1024 * 1024, // Disable swap
			CPUPeriod:  100000,
			CPUQuota:   50000,           // 0.5 CPU
			PidsLimit:  &[]int64{50}[0], // Limit processes
		},
		ReadonlyRootfs: true,
		NetworkMode:    "none",
		SecurityOpt: []string{
			"no-new-privileges",
		},
		CapDrop: []string{"ALL"},
	}

	containerResp, err := e.dockerClient.ContainerCreate(
		ctx, config, hostConfig, nil, nil, containerName,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create compilation container: %v", err)
	}

	e.trackContainer(containerResp.ID)

	if err := e.dockerClient.ContainerStart(ctx, containerResp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start compilation container: %v", err)
	}

	statusCh, errCh := e.dockerClient.ContainerWait(ctx, containerResp.ID, container.WaitConditionNotRunning)

	var compileOutput string
	select {
	case err := <-errCh:
		if err != nil {
			return "", fmt.Errorf("error waiting for compilation container: %v", err)
		}
	case status := <-statusCh:

		logReader, err := e.dockerClient.ContainerLogs(ctx, containerResp.ID, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		})
		if err != nil {
			return "", fmt.Errorf("failed to get compilation logs: %v", err)
		}
		defer logReader.Close()

		logs, err := io.ReadAll(logReader)
		if err != nil {
			return "", fmt.Errorf("failed to read compilation logs: %v", err)
		}

		if status.StatusCode != 0 {
			compileOutput = string(logs)
		}
	case <-time.After(30 * time.Second):
		// Timeout
		return "", fmt.Errorf("compilation timed out")
	}

	return compileOutput, nil
}

func (e *CodeExecutor) runTest(ctx context.Context, submissionID, submissionDir string, tc TestCase, timeLimit, memoryLimit int) (TestResult, error) {
	result := TestResult{
		TestCaseID: tc.ID,
		Passed:     false,
	}

	containerName := "run-" + submissionID + "-" + tc.ID

	inputPath := filepath.Join(submissionDir, "input.txt")
	if err := os.WriteFile(inputPath, []byte(tc.Input), 0644); err != nil {
		return result, fmt.Errorf("failed to write test input: %v", err)
	}

	if timeLimit <= 0 {
		timeLimit = 5
	}

	if memoryLimit <= 0 {
		memoryLimit = 128
	}

	config := &container.Config{
		Image: "gcc:latest",
		Cmd: []string{
			"sh", "-c",
			"cd /code && ./solution < input.txt",
		},
		Tty:          false,
		AttachStdout: true,
		AttachStderr: true,
	}

	hostConfig := &container.HostConfig{
		AutoRemove: false,
		Mounts: []mount.Mount{
			{
				Type:     mount.TypeBind,
				Source:   submissionDir,
				Target:   "/code",
				ReadOnly: false,
			},
		},
		Resources: container.Resources{
			Memory:     int64(memoryLimit) * 1024 * 1024,
			MemorySwap: int64(memoryLimit) * 1024 * 1024, // Disable swap
			CPUPeriod:  100000,
			CPUQuota:   50000,           // 0.5 CPU
			PidsLimit:  &[]int64{30}[0], // Limit processes
		},
		ReadonlyRootfs: true,
		NetworkMode:    "none",
		SecurityOpt: []string{
			"no-new-privileges",
			"seccomp=/code/seccomp.json",
		},
		CapDrop: []string{"ALL"},
	}

	containerResp, err := e.dockerClient.ContainerCreate(
		ctx, config, hostConfig, nil, nil, containerName,
	)
	if err != nil {
		return result, fmt.Errorf("failed to create execution container: %v", err)
	}

	e.trackContainer(containerResp.ID)

	if err := e.dockerClient.ContainerStart(ctx, containerResp.ID, container.StartOptions{}); err != nil {
		return result, fmt.Errorf("failed to start execution container: %v", err)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeLimit)*time.Second)
	defer cancel()

	statusCh, errCh := e.dockerClient.ContainerWait(ctx, containerResp.ID, container.WaitConditionNotRunning)

	var exitCode int64
	select {
	case err := <-errCh:
		if err != nil {
			return result, fmt.Errorf("error waiting for execution container: %v", err)
		}
	case status := <-statusCh:
		exitCode = status.StatusCode
	case <-timeoutCtx.Done():

		seconds := 5
		e.dockerClient.ContainerStop(ctx, containerResp.ID, container.StopOptions{Timeout: &seconds})
		result.Error = "Time limit exceeded"
		return result, nil
	}

	logReader, err := e.dockerClient.ContainerLogs(ctx, containerResp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return result, fmt.Errorf("failed to get execution logs: %v", err)
	}
	defer logReader.Close()

	logs, err := io.ReadAll(logReader)
	if err != nil {
		return result, fmt.Errorf("failed to read execution logs: %v", err)
	}

	programOutput := cleanOutput(string(logs))
	result.Output = programOutput

	if exitCode != 0 {
		result.Error = fmt.Sprintf("Runtime error: program exited with code %d", exitCode)
		return result, nil
	}

	expectedOutput := cleanOutput(tc.ExpectedOutput)
	result.Passed = programOutput == expectedOutput
	if !result.Passed {
		result.Error = fmt.Sprintf("Expected '%s' but got '%s'",
			formatForDisplay(expectedOutput),
			formatForDisplay(programOutput))
	}

	return result, nil
}

func (e *CodeExecutor) trackContainer(containerID string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.containers[containerID] = time.Now()
}

func (e *CodeExecutor) cleanupContainers() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		e.mu.Lock()
		now := time.Now()
		for id, created := range e.containers {
			if now.Sub(created) > 30*time.Minute {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				e.dockerClient.ContainerRemove(ctx, id, container.RemoveOptions{Force: true})
				cancel()
				delete(e.containers, id)
			}
		}
		e.mu.Unlock()
	}
}

type Server struct {
	challenges    map[string]Challenge
	challengesDir string
	executor      *CodeExecutor
}

func NewServer(challengesDir string, executor *CodeExecutor) (*Server, error) {
	if _, err := os.Stat(challengesDir); os.IsNotExist(err) {
		if err := os.MkdirAll(challengesDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create challenges directory: %v", err)
		}
		createSampleChallenge(challengesDir)
	}

	challenges := loadChallenges(challengesDir)
	log.Printf("Loaded %d challenges", len(challenges))

	return &Server{
		challenges:    challenges,
		challengesDir: challengesDir,
		executor:      executor,
	}, nil
}

func (s *Server) GetChallenges(req *struct{}, resp *map[string]Challenge) error {
	*resp = s.challenges
	return nil
}

func (s *Server) GetChallenge(id *string, resp *Challenge) error {
	challenge, ok := s.challenges[*id]
	if !ok {
		return fmt.Errorf("challenge not found: %s", *id)
	}
	*resp = challenge
	return nil
}

func (s *Server) SubmitSolution(req *SubmissionRequest, resp *SubmissionResponse) error {
	challenge, ok := s.challenges[req.ChallengeID]
	if !ok {
		return fmt.Errorf("challenge not found: %s", req.ChallengeID)
	}

	submissionID := uuid.New().String()

	execReq := ExecutionRequest{
		Code:         req.Code,
		Challenge:    challenge,
		SubmissionID: submissionID,
	}

	var execResp ExecutionResponse
	if err := s.executor.Execute(&execReq, &execResp); err != nil {
		return err
	}

	if execResp.Error != "" {
		resp.Success = false
		resp.Message = execResp.Error
		return nil
	}

	resp.TestResults = execResp.Results
	resp.Success = allTestsPassed(resp.TestResults)
	resp.Message = "Submission processed"

	return nil
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

func allTestsPassed(results []TestResult) bool {
	for _, result := range results {
		if !result.Passed {
			return false
		}
	}
	return len(results) > 0
}

// HTTP handler for the RPC server
type RPCHttpHandler struct {
	rpcServer *rpc.Server
}

func (h *RPCHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "CONNECT" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "405 must CONNECT\n")
		return
	}

	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Print("rpc hijacking error: ", err)
		return
	}

	io.WriteString(conn, "HTTP/1.0 200 Connected to RPC\r\n\r\n")
	h.rpcServer.ServeConn(conn)
}

func apiHandler(server *Server) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/challenges", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(server.challenges)
	})

	mux.HandleFunc("/api/challenge/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/api/challenge/"):]
		if challenge, ok := server.challenges[id]; ok {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(challenge)
		} else {
			http.NotFound(w, r)
		}
	})

	mux.HandleFunc("/api/submit", func(w http.ResponseWriter, r *http.Request) {
		var req SubmissionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var resp SubmissionResponse
		if err := server.SubmitSolution(&req, &resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	return mux
}

func corsMiddleware(next http.Handler) http.Handler {
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

var seccompProfile = `{
	"defaultAction": "SCMP_ACT_ERRNO",
	"architectures": [
		"SCMP_ARCH_X86_64",
		"SCMP_ARCH_X86",
		"SCMP_ARCH_AARCH64"
	],
	"syscalls": [
		{
			"names": [
				"access",
				"arch_prctl",
				"brk",
				"close",
				"execve",
				"exit_group",
				"fstat",
				"futex",
				"getcwd",
				"getegid",
				"geteuid",
				"getgid",
				"getpid",
				"getppid",
				"getuid",
				"lseek",
				"mmap",
				"mprotect",
				"munmap",
				"openat",
				"pread64",
				"read",
				"readlink",
				"rt_sigaction",
				"rt_sigprocmask",
				"set_robust_list",
				"set_tid_address",
				"stat",
				"statfs",
				"write"
			],
			"action": "SCMP_ACT_ALLOW"
		}
	]
}`

func main() {
	log.Println("Starting C code execution server...")

	workDir := "./executor-workdir"
	if err := os.MkdirAll(workDir, 0755); err != nil {
		log.Fatalf("Failed to create work directory: %v", err)
	}

	executor, err := NewCodeExecutor(workDir)
	if err != nil {
		log.Fatalf("Failed to create code executor: %v", err)
	}

	challengesDir := "./challenges"
	server, err := NewServer(challengesDir, executor)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	rpcServer := rpc.NewServer()
	if err := rpcServer.Register(server.executor); err != nil {
		log.Fatalf("Failed to register RPC executor: %v", err)
	}
	if err := rpcServer.Register(server); err != nil {
		log.Fatalf("Failed to register RPC server: %v", err)
	}

	rpcHandler := &RPCHttpHandler{rpcServer: rpcServer}

	apiMux := apiHandler(server)

	// Create main HTTP server
	httpMux := http.NewServeMux()
	httpMux.Handle("/rpc", rpcHandler)
	httpMux.Handle("/api/", corsMiddleware(apiMux))
	httpMux.Handle("/", http.FileServer(http.Dir("./frontend")))

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to start RPC listener: %v", err)
	}
	go rpcServer.Accept(listener)
	log.Println("RPC server started on port 9000")

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	log.Printf("HTTP server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, httpMux))
}
