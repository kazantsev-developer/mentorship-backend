package services

import (
	"github.com/kazantsev/mentorship-backend/internal/models"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
)

// BonusService manages bonus points for users
type BonusService struct {
	repo *repositories.BonusRepository
}

// NewBonusService returns a new BonusService instance
func NewBonusService(repo *repositories.BonusRepository) *BonusService {
	return &BonusService{repo: repo}
}

// GetBalance returns the current bonus balance for a user
func (s *BonusService) GetBalance(userID string) (int, error) {
	bal, err := s.repo.GetOrCreateBalance(userID)
	if err != nil {
		return 0, err
	}
	return bal.Balance, nil
}

// GetHistory returns the bonus transaction history for a user
func (s *BonusService) GetHistory(userID string) ([]models.BonusTransaction, error) {
	return s.repo.GetTransactions(userID)
}

// EarnAchievementBonus grants bonus points for unlocking an achievement
func (s *BonusService) EarnAchievementBonus(userID string, amount int, achievementID string) error {
	_, err := s.repo.AddTransaction(userID, models.BonusTypeAchievement, amount, "achievement reward", "achievement", achievementID)
	return err
}

// ConvertBonusToDiscount converts bonus points into a discount percentage
func (s *BonusService) ConvertBonusToDiscount(userID string, bonusAmount int) (int, error) {
	if bonusAmount <= 0 {
		return 0, nil
	}
	bal, err := s.repo.GetOrCreateBalance(userID)
	if err != nil {
		return 0, err
	}
	if bal.Balance < bonusAmount {
		return 0, nil
	}
	discount := bonusAmount / 100
	if discount > 15 {
		discount = 15
	}
	_, err = s.repo.AddTransaction(userID, models.BonusTypeDiscount, -bonusAmount, "converted to discount", "discount", "")
	if err != nil {
		return 0, err
	}
	return discount, nil
}
