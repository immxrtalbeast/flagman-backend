package supabase

import (
	"context"

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
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}
