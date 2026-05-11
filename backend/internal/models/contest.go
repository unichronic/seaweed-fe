package models

type Contest struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	EligibleTo            string `json:"eligible_to"`
	RegistrationStatus    string `json:"registration_status"`
	RegistrationStartTime int64  `json:"registration_start_time"`
	RegistrationEndTime   int64  `json:"registration_end_time"`
	StartTime             int64  `json:"start_time"`
	EndTime               int64  `json:"end_time"`
	Finalized             bool   `json:"finalized"`
}

type Problem struct {
	ID          string     `json:"id"`
	ContestID   string     `json:"contest_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Score       int        `json:"score"`
	TestCases   []TestCase `json:"test_cases,omitempty"`
}

type Ranking struct {
	ContestID         string `json:"contest_id"`
	UserID            string `json:"user_id"`
	Name              string `json:"name"`
	USN               string `json:"usn"`
	Department        string `json:"department"`
	Score             int    `json:"score"`
	Hidden            bool   `json:"hidden"`
	Disqualified      bool   `json:"disqualified"`
	Shortlisted       bool   `json:"shortlisted"`
	CorrectAttempts   int    `json:"correct_attempts"`
	IncorrectAttempts int    `json:"incorrect_attempts"`
	Rank              int    `json:"rank"`
}

type ContestRegistration struct {
	ContestID    string `json:"contest_id"`
	UserID       string `json:"user_id"`
	RegisteredAt int64  `json:"registered_at"`
}
