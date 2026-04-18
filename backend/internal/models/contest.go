package models

type Contest struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	EligibleTo            string `json:"eligible_to"`
	RegistrationStartTime int64  `json:"registration_start_time"`
	RegistrationEndTime   int64  `json:"registration_end_time"`
	StartTime             int64  `json:"start_time"`
	EndTime               int64  `json:"end_time"`
}

type Problem struct {
	ID          string `json:"id"`
	ContestID   string `json:"contest_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Score       int    `json:"score"`
}
