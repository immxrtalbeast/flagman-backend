package enterprise

import (
	"context"
	"fmt"
	"time"

	"github.com/immxrtalbeast/flagman-backend/internal/domain"
)

type EnterpriseInteractor struct {
	entrRepo domain.EnterpriseRepository
	usrRepo  domain.UserRepository
}

func NewEnterpriseInteractor(entrRepo domain.EnterpriseRepository, usrRepo domain.UserRepository) *EnterpriseInteractor {
	return &EnterpriseInteractor{entrRepo: entrRepo, usrRepo: usrRepo}
}

func (ent *EnterpriseInteractor) CreateEnterprise(userID uint, name string, description string) (*domain.Enterprise, error) {
	const op = "uc.enterprise.create"
	enterprise := domain.Enterprise{
		Name:        name,
		Description: description,
		CreatorID:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := ent.entrRepo.CreateEnterprise(&enterprise)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	user, err := ent.usrRepo.User(context.Background(), userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	err = ent.entrRepo.AddUserEnterprise(user.ID, enterprise.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &enterprise, nil
}

func (ent *EnterpriseInteractor) EnterpriseByID(id uint) (*domain.Enterprise, error) {
	const op = "uc.enterprise.byID"
	enterprise, err := ent.entrRepo.EnterpriseByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return enterprise, nil

}
