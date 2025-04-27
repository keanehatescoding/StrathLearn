package config

import (
	"os"
)

type Config struct {
	Port            string
	ChallengesDir   string
	UseDocker       bool
	DockerImage     string
	AuthEnabled     bool
	AuthProviderURL string
	AuthCallbackURL string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	challengesDir := os.Getenv("CHALLENGES_DIR")
	if challengesDir == "" {
		challengesDir = "./backend/challenges"
	}

	useDocker := os.Getenv("USE_DOCKER") != "false"

	dockerImage := os.Getenv("DOCKER_IMAGE")
	if dockerImage == "" {
		dockerImage = "strathlearn-code-runner:latest"
	}

	authEnabled := os.Getenv("AUTH_ENABLED") == "true"

	authProviderURL := os.Getenv("AUTH_PROVIDER_URL")
	authCallbackURL := os.Getenv("AUTH_CALLBACK_URL")

	return &Config{
		Port:            port,
		ChallengesDir:   challengesDir,
		UseDocker:       useDocker,
		DockerImage:     dockerImage,
		AuthEnabled:     authEnabled,
		AuthProviderURL: authProviderURL,
		AuthCallbackURL: authCallbackURL,
	}
}
