package runner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/google/uuid"

	"strathlearn/backend/models"
	"strathlearn/backend/utils"
)

var containerCleanupChannel = make(chan string, 100)

type DockerRunner struct {
	client *client.Client
}

func NewDockerRunner(client *client.Client) *DockerRunner {
	return &DockerRunner{client: client}
}

func (r *DockerRunner) IsDockerAvailable() bool {
	return r.client != nil
}

func StartContainerCleanupWorker(client *client.Client) {
	ctx := context.Background()
	for containerId := range containerCleanupChannel {
		time.Sleep(500 * time.Millisecond)
		err := client.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{Force: true})
		if err != nil && !strings.Contains(err.Error(), "No such container") {
			log.Printf("Error removing container %s: %v", containerId, err)
		}
	}
}

func scheduleContainerCleanup(containerId string) {
	containerCleanupChannel <- containerId
}

func (r *DockerRunner) RunTests(code string, challenge models.Challenge) []models.TestResult {
	log.Printf("Running tests for challenge: %s", challenge.ID)
	results := make([]models.TestResult, 0, len(challenge.TestCases))
	submissionID := uuid.New().String()

	codeDir := "/code/" + submissionID
	sourcePath := fmt.Sprintf("%s/solution.c", codeDir)

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

	setupResp, err := r.client.ContainerCreate(
		ctx, setupConfig, hostConfig, nil, nil, "setup-"+submissionID)
	if err != nil {
		log.Printf("Error creating setup container: %v", err)
		return []models.TestResult{{
			TestCaseID: "setup",
			Passed:     false,
			Error:      "System error: could not create setup container",
		}}
	}

	attachResp, err := r.client.ContainerAttach(
		ctx, setupResp.ID, types.ContainerAttachOptions{
			Stdin:  true,
			Stdout: true,
			Stderr: true,
			Stream: true,
		})
	if err != nil {
		log.Printf("Error attaching to setup container: %v", err)
		scheduleContainerCleanup(setupResp.ID)
		return []models.TestResult{{
			TestCaseID: "setup",
			Passed:     false,
			Error:      "System error: could not attach to setup container",
		}}
	}
	defer attachResp.Close()

	if err := r.client.ContainerStart(
		ctx, setupResp.ID, types.ContainerStartOptions{}); err != nil {
		log.Printf("Error starting setup container: %v", err)
		scheduleContainerCleanup(setupResp.ID)
		return []models.TestResult{{
			TestCaseID: "setup",
			Passed:     false,
			Error:      "System error: could not start setup container",
		}}
	}

	go func() {
		io.Copy(attachResp.Conn, strings.NewReader(code))
		attachResp.CloseWrite()
	}()

	statusCh, errCh := r.client.ContainerWait(
		ctx, setupResp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Printf("Error waiting for setup container: %v", err)
			scheduleContainerCleanup(setupResp.ID)
			return []models.TestResult{{
				TestCaseID: "setup",
				Passed:     false,
				Error:      "System error: setup container error",
			}}
		}
	case <-statusCh:
	}

	logs, err := r.client.ContainerLogs(
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

	compileResp, err := r.client.ContainerCreate(
		ctx, compileConfig, hostConfig, nil, nil, "compile-"+submissionID)
	if err != nil {
		log.Printf("Error creating compile container: %v", err)
		return []models.TestResult{{
			TestCaseID: "compile",
			Passed:     false,
			Error:      "System error: could not create compile container",
		}}
	}

	if err := r.client.ContainerStart(
		ctx, compileResp.ID, types.ContainerStartOptions{}); err != nil {
		log.Printf("Error starting compile container: %v", err)
		scheduleContainerCleanup(compileResp.ID)
		return []models.TestResult{{
			TestCaseID: "compile",
			Passed:     false,
			Error:      "System error: could not start compile container",
		}}
	}

	statusCh, errCh = r.client.ContainerWait(
		ctx, compileResp.ID, container.WaitConditionNotRunning)
	var compileExitCode int64
	select {
	case err := <-errCh:
		if err != nil {
			log.Printf("Error waiting for compile container: %v", err)
			scheduleContainerCleanup(compileResp.ID)
			return []models.TestResult{{
				TestCaseID: "compile",
				Passed:     false,
				Error:      "System error: compile container error",
			}}
		}
	case status := <-statusCh:
		compileExitCode = status.StatusCode
	}

	compileOutput, err := r.client.ContainerLogs(
		ctx, compileResp.ID, types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		})

	var compileErrorMsg string
	if err == nil {
		defer compileOutput.Close()
		var outputBuf bytes.Buffer
		_, err = utils.StdCopy(&outputBuf, &outputBuf, compileOutput)
		if err != nil {
			log.Printf("Error reading compile logs: %v", err)
		} else if compileExitCode != 0 {
			compileErrorMsg = outputBuf.String()
			scheduleContainerCleanup(compileResp.ID)
			return []models.TestResult{{
				TestCaseID: "compile",
				Passed:     false,
				Error:      "Compilation error: " + compileErrorMsg,
			}}
		}
	}

	scheduleContainerCleanup(compileResp.ID)

	for _, tc := range challenge.TestCases {
		result := models.TestResult{
			TestCaseID: tc.ID,
			Passed:     false,
		}

		inputPath := fmt.Sprintf("%s/input.txt", codeDir)
		inputConfig := &container.Config{
			Image: "strathlearn-code-runner:latest",
			Cmd: []string{
				"sh", "-c",
				fmt.Sprintf("echo -n %s > %s",
					utils.Shellescape(tc.Input),
					inputPath),
			},
			Tty: false,
		}

		inputResp, err := r.client.ContainerCreate(
			ctx, inputConfig, hostConfig, nil, nil, "input-"+submissionID+"-"+tc.ID)
		if err != nil {
			log.Printf("Error creating input container: %v", err)
			result.Error = "System error: failed to prepare input"
			results = append(results, result)
			continue
		}

		if err := r.client.ContainerStart(
			ctx, inputResp.ID, types.ContainerStartOptions{}); err != nil {
			log.Printf("Error starting input container: %v", err)
			scheduleContainerCleanup(inputResp.ID)
			result.Error = "System error: failed to prepare input"
			results = append(results, result)
			continue
		}

		r.client.ContainerWait(
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

		runResp, err := r.client.ContainerCreate(
			ctx, runConfig, runHostConfig, nil, nil, "run-"+submissionID+"-"+tc.ID)
		if err != nil {
			log.Printf("Error creating run container: %v", err)
			result.Error = "System error: could not create run container"
			results = append(results, result)
			continue
		}

		if err := r.client.ContainerStart(
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
			r.client.ContainerStop(ctx, runResp.ID, container.StopOptions{})
			result.Error = "Time limit exceeded"
		case err := <-func() <-chan error {
			ch := make(chan error, 1)
			go func() {
				statusCh, errCh := r.client.ContainerWait(
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

		output, err := r.client.ContainerLogs(
			ctx, runResp.ID, types.ContainerLogsOptions{
				ShowStdout: true,
				ShowStderr: true,
			})
		if err == nil {
			defer output.Close()
			var outputBuf bytes.Buffer
			_, err = utils.StdCopy(&outputBuf, &outputBuf, output)
			if err != nil {
				log.Printf("Error reading output: %v", err)
				result.Error = "System error: failed to read program output"
			} else {
				programOutput := utils.CleanOutput(outputBuf.String())
				result.Output = programOutput

				if result.Error == "" {
					expectedOutput := utils.CleanOutput(tc.ExpectedOutput)

					log.Printf("Expected output: '%s'", expectedOutput)
					log.Printf("Actual output: '%s'", programOutput)

					result.Passed = programOutput == expectedOutput
					if !result.Passed {
						result.Error = fmt.Sprintf("Expected '%s' but got '%s'",
							utils.FormatForDisplay(expectedOutput),
							utils.FormatForDisplay(programOutput))
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
