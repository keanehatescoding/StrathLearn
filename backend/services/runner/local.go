package runner

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"strathlearn/backend/models"
	"strathlearn/backend/utils"
)

type LocalRunner struct{}

func NewLocalRunner() *LocalRunner {
	return &LocalRunner{}
}

func (r *LocalRunner) IsDockerAvailable() bool {
	return false
}

func (r *LocalRunner) RunTests(code string, challenge models.Challenge) []models.TestResult {
	log.Printf("Running tests locally for challenge: %s", challenge.ID)
	results := make([]models.TestResult, 0, len(challenge.TestCases))

	tempDir, err := os.MkdirTemp("", "challenge-")
	if err != nil {
		log.Printf("Failed to create temp dir: %v", err)
		return []models.TestResult{{
			TestCaseID: "compile",
			Passed:     false,
			Error:      "System error: could not create temporary directory",
		}}
	}
	defer os.RemoveAll(tempDir)

	submissionID := uuid.New().String()

	sourcePath := filepath.Join(tempDir, fmt.Sprintf("solution-%s.c", submissionID))
	if err := os.WriteFile(sourcePath, []byte(code), 0644); err != nil {
		return []models.TestResult{{
			TestCaseID: "compile",
			Passed:     false,
			Error:      "System error: could not write source code",
		}}
	}

	execPath := filepath.Join(tempDir, fmt.Sprintf("solution-%s", submissionID))
	cmd := exec.Command("gcc", sourcePath, "-o", execPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return []models.TestResult{{
			TestCaseID: "compile",
			Passed:     false,
			Error:      "Compilation error: " + string(output),
		}}
	}

	for _, tc := range challenge.TestCases {
		result := models.TestResult{
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
				programOutput := utils.CleanOutput(string(output))
				result.Output = programOutput

				expectedOutput := utils.CleanOutput(tc.ExpectedOutput)

				result.Passed = programOutput == expectedOutput
				if !result.Passed {
					result.Error = fmt.Sprintf("Expected '%s' but got '%s'",
						utils.FormatForDisplay(expectedOutput),
						utils.FormatForDisplay(programOutput))
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
