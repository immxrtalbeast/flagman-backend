package enterprise

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/immxrtalbeast/flagman-backend/internal/domain"
)

type EnterpriseInteractor struct {
	entrRepo  domain.EnterpriseRepository
	usrRepo   domain.UserRepository
	notifRepo domain.NotificationRepository
}

func NewEnterpriseInteractor(entrRepo domain.EnterpriseRepository, usrRepo domain.UserRepository, notifRepo domain.NotificationRepository) *EnterpriseInteractor {
	return &EnterpriseInteractor{entrRepo: entrRepo, usrRepo: usrRepo, notifRepo: notifRepo}
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

func (ent *EnterpriseInteractor) GetEnterprisesByUserID(userID uint) ([]domain.Enterprise, error) {
	const op = "uc.enterprise.byUserID"
	enterprises, err := ent.entrRepo.GetEnterprisesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return enterprises, nil
}

func (ent *EnterpriseInteractor) InviteUser(senderID uint, userEmail string, enterpriseID uint, enterpriseName string) (*domain.Invitation, error) {
	const op = "uc.enterprise.invite"
	enterpise, err := ent.entrRepo.EnterpriseByID(enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if enterpise.CreatorID != senderID {
		return nil, errors.New("недостаточно прав для отправки приглашений")
	}
	receiver, err := ent.usrRepo.FindByEmail(context.Background(), userEmail)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	invite := domain.Invitation{
		Email:          userEmail,
		EnterpriseID:   enterpriseID,
		CreatedBy:      senderID,
		EnterpriseName: enterpriseName,
		ReceiverID:     receiver.ID,
	}
	if err = ent.notifRepo.CreateInvitation(context.Background(), invite); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &invite, nil
}
