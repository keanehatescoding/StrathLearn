package test

import (
	"fmt"

	"strathlearn/backend/models"
	"strathlearn/backend/services/runner"
)

func main() {
	// Create a sample challenge
	challenge := models.Challenge{
		ID:          "test-challenge",
		Title:       "Test Challenge",
		TimeLimit:   2,
		MemoryLimit: 128,
		TestCases: []models.TestCase{
			{
				ID:             "test-case-1",
				Input:          "5",
				ExpectedOutput: "25",
				Hidden:         false,
			},
			{
				ID:             "test-case-2",
				Input:          "10",
				ExpectedOutput: "100",
				Hidden:         false,
			},
		},
	}

	// Sample submissions
	testSubmissions := []struct {
		name string
		code string
	}{
		{
			name: "Correct solution",
			code: `
				#include <stdio.h>
				
				int main() {
					int n;
					scanf("%d", &n);
					printf("%d", n * n);
					return 0;
				}
			`,
		},
		{
			name: "Wrong solution",
			code: `
				#include <stdio.h>
				
				int main() {
					int n;
					scanf("%d", &n);
					printf("%d", n + n);
					return 0;
				}
			`,
		},
		{
			name: "Compilation error",
			code: `
				#include <stdio.h>
				
				int main() {
					int n
					scanf("%d", &n);
					printf("%d", n * n);
					return 0;
				}
			`,
		},
		{
			name: "Timeout",
			code: `
				#include <stdio.h>
				
				int main() {
					int n;
					scanf("%d", &n);
					while(1) { }
					return 0;
				}
			`,
		},
	}

	// Create runner
	judge0Runner := runner.NewJudge0Runner()

	// Test each submission
	for _, submission := range testSubmissions {
		fmt.Printf("\n=== Testing: %s ===\n", submission.name)
		results := judge0Runner.RunTests(submission.code, challenge)

		for _, result := range results {
			fmt.Printf("Test case %s: %v\n", result.TestCaseID, result.Passed)
			if !result.Passed {
				fmt.Printf("  Error: %s\n", result.Error)
				fmt.Printf("  Output: %s\n", result.Output)
			}
		}
	}
}
