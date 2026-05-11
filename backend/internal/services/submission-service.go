package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AniketSrivastava1/recruit/backend/internal/models"
	"github.com/AniketSrivastava1/recruit/backend/internal/s3"
	"github.com/AniketSrivastava1/recruit/backend/internal/stores"
	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
)

type SubmissionService struct {
	store     *stores.SubmissionStore
	s3Client  *s3.S3Client
	probStore *stores.ContestStore
}

func NewSubmissionService(store *stores.SubmissionStore, s3Client *s3.S3Client, probStore *stores.ContestStore) *SubmissionService {
	return &SubmissionService{
		store:     store,
		s3Client:  s3Client,
		probStore: probStore,
	}
}

func (s *SubmissionService) Submit(ctx context.Context, sub *models.Submission, code string) (string, error) {
	if err := s.ensureCanSubmit(ctx, sub.ContestID, sub.UserID); err != nil {
		return "", err
	}
	if _, err := s.probStore.GetProblem(ctx, sub.ContestID, sub.ProblemID, true); err != nil {
		return "", err
	}
	sub.ID = ulid.Make().String()
	sub.CreatedAt = time.Now().UnixMilli()
	sub.Status = "pending"
	sub.S3Key = "submissions/" + sub.ContestID + "/" + sub.UserID + "/" + sub.ID + extensionForLanguage(sub.Language)

	if err := s.s3Client.UploadData(ctx, sub.S3Key, io.Reader(bytes.NewReader([]byte(code)))); err != nil {
		return "", err
	}

	if err := s.store.CreateSubmission(ctx, sub); err != nil {
		return "", err
	}

	go s.judge(sub, code)

	return sub.ID, nil
}

func (s *SubmissionService) GetStatus(ctx context.Context, submissionId, userId string) (*models.Submission, error) {
	return s.store.GetSubmission(ctx, submissionId, userId)
}

func (s *SubmissionService) GetDetails(ctx context.Context, submissionId, userId string) (map[string]any, error) {
	sub, err := s.store.GetSubmission(ctx, submissionId, userId)
	if err != nil {
		return nil, err
	}
	results, err := s.store.GetTestCaseResults(ctx, submissionId)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"submission":        sub,
		"test_case_results": results,
	}, nil
}

func (s *SubmissionService) List(ctx context.Context, userId, problemId string) ([]models.Submission, error) {
	return s.store.ListSubmissions(ctx, userId, problemId)
}

func (s *SubmissionService) judge(sub *models.Submission, code string) {
	ctx := context.Background()
	problem, err := s.probStore.GetProblem(ctx, sub.ContestID, sub.ProblemID, true)
	if err != nil {
		fmt.Printf("Failed to load problem for submission %s: %v\n", sub.ID, err)
		_ = s.store.UpdateStatus(ctx, sub.ID, "failed_to_process")
		return
	}

	status, results, err := s.runJudge0(ctx, sub, problem, code)
	if err != nil {
		fmt.Printf("Failed to judge submission %s: %v\n", sub.ID, err)
		status = "failed_to_process"
	}

	if err := s.store.ApplyVerdict(ctx, sub, status, problem.Score, results); err != nil {
		fmt.Printf("Failed to apply verdict: %v\n", err)
		return
	}
}

func (s *SubmissionService) ensureCanSubmit(ctx context.Context, contestId, userId string) error {
	contest, err := s.probStore.GetContest(ctx, contestId)
	if err != nil {
		return err
	}
	if contest.Finalized {
		return echo.NewHTTPError(http.StatusForbidden, "Contest is finalized")
	}
	registered, err := s.probStore.IsRegistered(ctx, contestId, userId)
	if err != nil {
		return err
	}
	if !registered {
		return echo.NewHTTPError(http.StatusForbidden, "Register for the contest before submitting")
	}
	now := time.Now().UnixMilli()
	if now < contest.StartTime || now > contest.EndTime {
		return echo.NewHTTPError(http.StatusForbidden, "Submissions are allowed only during the contest window")
	}
	return nil
}

func (s *SubmissionService) runJudge0(ctx context.Context, sub *models.Submission, problem *models.Problem, code string) (string, []models.TestCaseResult, error) {
	baseURL := os.Getenv("JUDGE0_BASE_URL")
	if baseURL == "" {
		return "failed_to_process", nil, fmt.Errorf("JUDGE0_BASE_URL not set")
	}
	if len(problem.TestCases) == 0 {
		return "failed_to_process", nil, fmt.Errorf("problem has no test cases")
	}
	languageID, ok := judge0LanguageID(sub.Language)
	if !ok {
		return "failed_to_process", nil, fmt.Errorf("unsupported language: %s", sub.Language)
	}

	results := make([]models.TestCaseResult, 0, len(problem.TestCases))
	finalStatus := "accepted"
	for i, tc := range problem.TestCases {
		testID := tc.ID
		if testID == "" {
			testID = fmt.Sprintf("case-%d", i+1)
		}
		status, runtime, memory, err := callJudge0(ctx, baseURL, languageID, tc, code)
		if err != nil {
			return "failed_to_process", results, err
		}
		results = append(results, models.TestCaseResult{
			SubmissionID: sub.ID,
			TestCaseID:   testID,
			Status:       status,
			RuntimeMS:    runtime,
			MemoryKB:     memory,
			CreatedAt:    time.Now().UnixMilli(),
		})
		if status != "pass" && finalStatus == "accepted" {
			finalStatus = status
		}
	}
	return finalStatus, results, nil
}

type judge0Request struct {
	SourceCode     string `json:"source_code"`
	LanguageID     int    `json:"language_id"`
	Stdin          string `json:"stdin"`
	ExpectedOutput string `json:"expected_output"`
}

type judge0Response struct {
	Stdout        string `json:"stdout"`
	Stderr        string `json:"stderr"`
	CompileOutput string `json:"compile_output"`
	Time          string `json:"time"`
	Memory        int64  `json:"memory"`
	Status        struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	} `json:"status"`
}

func callJudge0(ctx context.Context, baseURL string, languageID int, tc models.TestCase, code string) (string, int64, int64, error) {
	reqBody := judge0Request{
		SourceCode:     code,
		LanguageID:     languageID,
		Stdin:          tc.Stdin,
		ExpectedOutput: tc.ExpectedOutput,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return "failed_to_process", 0, 0, err
	}
	judgeCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(judgeCtx, http.MethodPost, strings.TrimRight(baseURL, "/")+"/submissions?base64_encoded=false&wait=true", bytes.NewReader(body))
	if err != nil {
		return "failed_to_process", 0, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "failed_to_process", 0, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "failed_to_process", 0, 0, fmt.Errorf("judge0 returned status %d", resp.StatusCode)
	}
	var judgeResp judge0Response
	if err := json.NewDecoder(resp.Body).Decode(&judgeResp); err != nil {
		return "failed_to_process", 0, 0, err
	}
	return mapJudge0Status(judgeResp.Status.ID), parseRuntimeMS(judgeResp.Time), judgeResp.Memory, nil
}

func judge0LanguageID(language string) (int, bool) {
	switch language {
	case "python":
		return 71, true
	case "cpp":
		return 54, true
	case "java":
		return 62, true
	default:
		return 0, false
	}
}

func mapJudge0Status(id int) string {
	switch id {
	case 3:
		return "pass"
	case 4:
		return "wrong_answer"
	case 5:
		return "tle"
	case 6, 7, 8, 9, 10, 11, 12, 13, 14:
		return "rte"
	default:
		return "wrong_answer"
	}
}

func parseRuntimeMS(value string) int64 {
	var seconds float64
	if _, err := fmt.Sscanf(value, "%f", &seconds); err != nil {
		return 0
	}
	return int64(seconds * 1000)
}

func extensionForLanguage(language string) string {
	switch language {
	case "python":
		return ".py"
	case "cpp":
		return ".cpp"
	case "java":
		return ".java"
	default:
		return ".txt"
	}
}
