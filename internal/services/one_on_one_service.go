package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/gorm"
)

// OneOnOneService manages 1:1 session requests and approvals
type OneOnOneService struct {
	db           *gorm.DB
	bonusService *BonusService
}

// NewOneOnOneService returns a new OneOnOneService instance
func NewOneOnOneService(db *gorm.DB, bonusService *BonusService) *OneOnOneService {
	return &OneOnOneService{db: db, bonusService: bonusService}
}

// CreateRequest submits a new 1:1 session request for a student
func (s *OneOnOneService) CreateRequest(studentID string) (*models.OneOnOneRequest, error) {
	req := models.OneOnOneRequest{
		ID:        uuid.New().String(),
		StudentID: studentID,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.db.Create(&req).Error
	return &req, err
}

// GetAllRequests returns all 1:1 requests ordered by most recent
func (s *OneOnOneService) GetAllRequests() ([]models.OneOnOneRequest, error) {
	var requests []models.OneOnOneRequest
	err := s.db.Order("created_at desc").Find(&requests).Error
	return requests, err
}

// ApproveRequest approves a pending 1:1 request after checking the student's bonus balance
func (s *OneOnOneService) ApproveRequest(requestID string, adminID string) error {
	var req models.OneOnOneRequest
	if err := s.db.First(&req, "id = ?", requestID).Error; err != nil {
		return err
	}
	if req.Status != "pending" {
		return errors.New("request already processed")
	}
	balance, err := s.bonusService.GetBalance(req.StudentID)
	if err != nil {
		return err
	}
	if balance < 1000 {
		return errors.New("insufficient bonus balance (need 1000)")
	}
	_, err = s.bonusService.repo.AddTransaction(req.StudentID, models.BonusTypeOneOnOneSpend, -1000, "1x1 session", "one_on_one", requestID)
	if err != nil {
		return err
	}
	req.Status = "approved"
	req.ProcessedBy = &adminID
	now := time.Now()
	req.ProcessedAt = &now
	return s.db.Save(&req).Error
}

// RejectRequest rejects a pending 1:1 request
func (s *OneOnOneService) RejectRequest(requestID string, adminID string) error {
	var req models.OneOnOneRequest
	if err := s.db.First(&req, "id = ?", requestID).Error; err != nil {
		return err
	}
	if req.Status != "pending" {
		return errors.New("request already processed")
	}
	req.Status = "rejected"
	req.ProcessedBy = &adminID
	now := time.Now()
	req.ProcessedAt = &now
	return s.db.Save(&req).Error
}
