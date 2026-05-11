package models

type CreateUserRequest struct {
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	USN          string `json:"usn" validate:"required"`
	MobileNumber string `json:"mobile_number"`
	JoiningYear  int    `json:"joining_year" validate:"required,gt=0"`
	Department   string `json:"department" validate:"required"`
}

type UpdateUserProfileRequest struct {
	Name         string `json:"name" validate:"omitempty"`
	MobileNumber string `json:"mobile_number" validate:"omitempty"`
	Department   string `json:"department" validate:"omitempty"`
	JoiningYear  int    `json:"joining_year" validate:"omitempty,gt=0"`
}

type CreateContestRequest struct {
	Name                  string `json:"name" validate:"required"`
	Description           string `json:"description"`
	EligibleTo            string `json:"eligible_to"`
	RegistrationStatus    string `json:"registration_status" validate:"omitempty,oneof=open closed invite-only"`
	RegistrationStartTime int64  `json:"registration_start_time" validate:"required,gt=0"`
	RegistrationEndTime   int64  `json:"registration_end_time" validate:"required,gt=0"`
	StartTime             int64  `json:"start_time" validate:"required,gt=0"`
	EndTime               int64  `json:"end_time" validate:"required,gt=0"`
}

type ContestRegistrationRequest struct {
	Action string `json:"action" validate:"required,oneof=register unregister"`
}

type TestCase struct {
	ID             string `json:"id"`
	Stdin          string `json:"stdin"`
	ExpectedOutput string `json:"expected_output"`
}

type ProblemRequest struct {
	Name        string     `json:"name" validate:"required"`
	Description string     `json:"description" validate:"required"`
	Score       int        `json:"score" validate:"required,gt=0"`
	TestCases   []TestCase `json:"test_cases" validate:"required,min=1,dive"`
}

type SubmissionRequest struct {
	ContestID string `json:"contest_id" validate:"required"`
	ProblemID string `json:"problem_id" validate:"required"`
	Language  string `json:"language" validate:"required,oneof=python cpp java"`
	Code      string `json:"code" validate:"required"`
}

type LeaderboardUpdateRequest struct {
	Hidden       *bool `json:"hidden"`
	Disqualified *bool `json:"disqualified"`
	Shortlisted  *bool `json:"shortlisted"`
}

type ContestResponse struct {
	Contest
	Registered bool `json:"registered"`
}

type ProblemSummaryResponse struct {
	ID        string `json:"id"`
	ContestID string `json:"contest_id"`
	Name      string `json:"name"`
	Score     int    `json:"score"`
}

type ProblemDetailResponse struct {
	ID          string `json:"id"`
	ContestID   string `json:"contest_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Score       int    `json:"score"`
}
