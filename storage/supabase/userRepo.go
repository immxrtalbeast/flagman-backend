package supabase

import (
	"context"
	"fmt"

	"github.com/immxrtalbeast/flagman-backend/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) (uint, error) {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.ID, nil
}
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) User(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Enterprises").Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *UserRepository) Users(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) GetUsersByEnterpriseID(enterpriseID string) ([]domain.User, error) {
	var users []domain.User
	err := r.db.
		Joins("JOIN user_enterprises ON user_enterprises.user_id = users.id").
		Where("user_enterprises.enterprise_id = ?", enterpriseID).
		Find(&users).
		Error
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователей: %v", err)
	}

	return users, err
}
