package supabase

import (
	"github.com/immxrtalbeast/flagman-backend/internal/domain"
	"gorm.io/gorm"
)

type EnterpriseRepository struct {
	db *gorm.DB
}

func NewEnterpriseRepository(db *gorm.DB) *EnterpriseRepository {
	return &EnterpriseRepository{db: db}
}

func (r *EnterpriseRepository) CreateEnterprise(enterprise *domain.Enterprise) error {
	const op = "repo.enterprise.create"
	result := r.db.Create(enterprise)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *EnterpriseRepository) AddUserEnterprise(userID uint, enterpriseID uint) error {
	// Найти организацию и пользователя
	var enterprise domain.Enterprise
	var user domain.User

	if err := r.db.First(&enterprise, enterpriseID).Error; err != nil {
		return err
	}

	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}
	return r.db.Model(&enterprise).Association("Users").Append(&user)
}

func (r *EnterpriseRepository) EnterpriseByID(enterpriseID uint) (*domain.Enterprise, error) {
	const op = "repo.enterprise.byID"
	var enterprise domain.Enterprise
	if err := r.db.Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Omit("PassHash", "PhoneNumber", "Enterprises", "Roles", "SentDocuments", "ReceivedDocuments")
	}).First(&enterprise, enterpriseID).Error; err != nil {
		return nil, err
	}
	return &enterprise, nil
}
func (r *EnterpriseRepository) GetEnterprisesByUserID(userID uint) ([]domain.Enterprise, error) {
	var enterprises []domain.Enterprise
	err := r.db.
		Joins("JOIN user_enterprises ON user_enterprises.enterprise_id = enterprises.id").
		Where("user_enterprises.user_id = ?", userID).
		Find(&enterprises).
		Error
	return enterprises, err
}
