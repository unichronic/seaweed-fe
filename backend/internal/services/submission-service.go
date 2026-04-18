package services

import (
	"bytes"
	"context"
	"io"
	"log"
	"time"

	"github.com/AniketSrivastava1/recruit/backend/internal/models"
	"github.com/AniketSrivastava1/recruit/backend/internal/s3"
	"github.com/AniketSrivastava1/recruit/backend/internal/stores"
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
	sub.ID = ulid.Make().String()
	sub.CreatedAt = time.Now().UnixMilli()
	sub.Status = "pending"
	sub.S3Key = "submissions/" + sub.ContestID + "/" + sub.UserID + "/" + sub.ID + ".txt"

	// 1. Upload to S3
	if err := s.s3Client.UploadData(ctx, sub.S3Key, io.Reader(bytes.NewReader([]byte(code)))); err != nil {
		return "", err
	}

	// 2. Save to DB
	if err := s.store.CreateSubmission(ctx, sub); err != nil {
		return "", err
	}

	// 3. Async Judging
	go s.judge(sub)

	return sub.ID, nil
}

func (s *SubmissionService) judge(sub *models.Submission) {
	// Simulate Judge0 call for now
	time.Sleep(5 * time.Second)

	// Update status to accepted (for demo)
	status := "accepted"
	ctx := context.Background()

	if err := s.store.UpdateStatus(ctx, sub.ID, status); err != nil {
		log.Printf("Failed to update status: %v", err)
		return
	}

	if status == "accepted" {
		points := 10 // mock

		if err := s.store.UpdateRanking(ctx, sub.ContestID, sub.UserID, points); err != nil {
			log.Printf("Failed to update ranking: %v", err)
		}

		if err := s.store.RefreshRankingMV(ctx); err != nil {
			log.Printf("Failed to refresh ranking MV: %v", err)
		}
	}
}
