package runner

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CodeExecutionRequest struct {
	Code        string `json:"code" binding:"required"`
	LanguageID  int    `json:"language_id" binding:"required"`
	Input       string `json:"input"`
	TimeLimit   int    `json:"time_limit"`
	MemoryLimit int    `json:"memory_limit"`
}

type CodeExecutionResponse struct {
	Output        string  `json:"output"`
	Error         string  `json:"error"`
	ExecutionTime float64 `json:"execution_time"`
	Memory        int     `json:"memory"`
	Status        string  `json:"status"`
	CompileOutput string  `json:"compile_output"`
}

func ExecuteCode(c *gin.Context) {
	var req CodeExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default limits if not provided
	if req.TimeLimit == 0 {
		req.TimeLimit = 5 // 5 seconds default
	}
	if req.MemoryLimit == 0 {
		req.MemoryLimit = 128 // 128MB default
	}

	// Create Judge0 runner
	judge0Runner := NewJudge0Runner()

	// Submit code for execution
	token, err := judge0Runner.SubmitCodeExecution(req.Code, req.LanguageID, req.Input, req.TimeLimit, req.MemoryLimit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit code: " + err.Error()})
		return
	}

	// Wait for result
	result, err := judge0Runner.WaitForExecutionResult(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get execution result: " + err.Error()})
		return
	}

	response := CodeExecutionResponse{
		Output:        result.Stdout,
		Error:         result.Stderr,
		ExecutionTime: result.ExecutionTime,
		Memory:        result.Memory,
		Status:        result.Status.Description,
		CompileOutput: result.CompileOutput,
	}

	c.JSON(http.StatusOK, response)
}
