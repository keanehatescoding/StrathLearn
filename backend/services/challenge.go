package services

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"strathlearn/backend/models"
)

func LoadChallenges(dir string) map[string]models.Challenge {
	challenges := make(map[string]models.Challenge)
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

			var challenge models.Challenge
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
		CreateSampleChallenge(dir)
		return LoadChallenges(dir)
	}

	return challenges
}

func CreateSampleChallenge(dir string) {
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
