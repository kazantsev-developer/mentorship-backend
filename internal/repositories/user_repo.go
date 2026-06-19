package repositories

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/gorm"
)

// UserRepository handles database operations for user accounts and roles
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository returns a new UserRepository instance
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user along with their roles within a transaction
func (r *UserRepository) Create(user *models.User, roles []models.Role) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if user.ID == "" {
			user.ID = uuid.New().String()
		}
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		for _, role := range roles {
			userRole := models.UserRole{
				UserID: user.ID,
				Role:   role,
			}
			if err := tx.Create(&userRole).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// FindByLogin looks up a non-deleted user by login and returns their roles
func (r *UserRepository) FindByLogin(login string) (*models.User, []models.Role, error) {
	var user models.User
	err := r.db.Where("login = ? AND is_deleted = ?", login, false).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	var userRoles []models.UserRole
	if err := r.db.Where("user_id = ?", user.ID).Find(&userRoles).Error; err != nil {
		return nil, nil, err
	}

	roles := make([]models.Role, len(userRoles))
	for i, ur := range userRoles {
		roles[i] = ur.Role
	}
	return &user, roles, nil
}

// FindByID looks up a non-deleted user by ID
func (r *UserRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update saves changes to an existing user
func (r *UserRepository) Update(user *models.User) error {
	user.UpdatedAt = time.Now()
	return r.db.Save(user).Error
}

// SoftDelete marks a user as deleted
func (r *UserRepository) SoftDelete(id string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}

// GetDB returns the underlying database instance
func (r *UserRepository) GetDB() *gorm.DB {
	return r.db
}
