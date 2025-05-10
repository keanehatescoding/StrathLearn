package runner

import (
	"strathlearn/backend/models"

	"github.com/docker/docker/client"
)

type CodeRunner interface {
	RunTests(code string, challenge models.Challenge) []models.TestResult
	IsDockerAvailable() bool
}

func NewRunner(dockerClient *client.Client) CodeRunner {
	if dockerClient != nil {
		return NewDockerRunner(dockerClient)
	}
	return NewLocalRunner()
}
