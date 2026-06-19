package repositories

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/gorm"
)

// BonusRepository handles database operations for bonus balances and transactions
type BonusRepository struct {
	db *gorm.DB
}

// NewBonusRepository returns a new BonusRepository instance
func NewBonusRepository(db *gorm.DB) *BonusRepository {
	return &BonusRepository{db: db}
}

// GetOrCreateBalance returns the bonus balance for a user, creating one if not found
func (r *BonusRepository) GetOrCreateBalance(userID string) (*models.BonusBalance, error) {
	var balance models.BonusBalance
	err := r.db.Where("user_id = ?", userID).First(&balance).Error
	if err == nil {
		return &balance, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		balance = models.BonusBalance{
			ID:      uuid.New().String(),
			UserID:  userID,
			Balance: 0,
		}
		if err := r.db.Create(&balance).Error; err != nil {
			return nil, err
		}
		return &balance, nil
	}
	return nil, err
}

// AddTransaction records a bonus transaction and updates the user balance atomically
func (r *BonusRepository) AddTransaction(userID string, txType models.BonusTransactionType, amount int, reason, sourceType, sourceID string) (*models.BonusBalance, error) {
	var balance models.BonusBalance
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var b models.BonusBalance
		err := tx.Where("user_id = ?", userID).First(&b).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			b = models.BonusBalance{
				ID:      uuid.New().String(),
				UserID:  userID,
				Balance: 0,
			}
			if err := tx.Create(&b).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		if amount < 0 && b.Balance+amount < 0 {
			return errors.New("insufficient balance")
		}
		b.Balance += amount
		if err := tx.Save(&b).Error; err != nil {
			return err
		}
		txRecord := models.BonusTransaction{
			ID:         uuid.New().String(),
			UserID:     userID,
			Type:       txType,
			Amount:     amount,
			Reason:     reason,
			SourceType: sourceType,
			SourceID:   sourceID,
			CreatedAt:  time.Now(),
		}
		if err := tx.Create(&txRecord).Error; err != nil {
			return err
		}
		balance = b
		return nil
	})
	return &balance, err
}

// GetTransactions returns all bonus transactions for a user ordered by most recent
func (r *BonusRepository) GetTransactions(userID string) ([]models.BonusTransaction, error) {
	var txns []models.BonusTransaction
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&txns).Error
	return txns, err
}
