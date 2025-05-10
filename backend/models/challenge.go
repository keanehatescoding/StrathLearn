package models

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
