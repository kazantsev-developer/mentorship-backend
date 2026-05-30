package services

import (
	"github.com/kazantsev/mentorship-backend/internal/models"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
)

type BonusService struct {
	repo *repositories.BonusRepository
}

func NewBonusService(repo *repositories.BonusRepository) *BonusService {
	return &BonusService{repo: repo}
}

func (s *BonusService) GetBalance(userID string) (int, error) {
	bal, err := s.repo.GetOrCreateBalance(userID)
	if err != nil {
		return 0, err
	}
	return bal.Balance, nil
}

func (s *BonusService) GetHistory(userID string) ([]models.BonusTransaction, error) {
	return s.repo.GetTransactions(userID)
}

func (s *BonusService) EarnAchievementBonus(userID string, amount int, achievementID string) error {
	_, err := s.repo.AddTransaction(userID, models.BonusTypeAchievement, amount, "achievement reward", "achievement", achievementID)
	return err
}

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

func (s *BonusService) SpendForOneOnOne(userID, requestID string) (int, error) {
	bal, err := s.repo.AddTransaction(userID, models.BonusTypeOneOnOneSpend, -1000, "1x1 session", "one_on_one", requestID)
	if err != nil {
		return 0, err
	}
	return bal.Balance, nil
}
