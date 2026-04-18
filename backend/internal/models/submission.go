package models

type Submission struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	ContestID string `json:"contest_id"`
	ProblemID string `json:"problem_id"`
	Language  string `json:"language"`
	S3Key     string `json:"s3_key"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

type TestCaseResult struct {
	ID           string `json:"id"`
	SubmissionID string `json:"submission_id"`
	TestCaseID   string `json:"test_case_id"`
	Status       string `json:"status"`
	RuntimeMS    int64  `json:"runtime_ms"`
	MemoryKB     int64  `json:"memory_kb"`
	CreatedAt    int64  `json:"created_at"`
}
