package notifications

import (
	"context"
	"fmt"

	"github.com/immxrtalbeast/flagman-backend/internal/domain"
)

type NotificationInteractor struct {
	notifRepo domain.NotificationRepository
	entrRepo  domain.EnterpriseRepository
}

func NewNotificationInteractor(notifRepo domain.NotificationRepository, entrRepo domain.EnterpriseRepository) *NotificationInteractor {
	return &NotificationInteractor{notifRepo: notifRepo, entrRepo: entrRepo}
}

func (ni *NotificationInteractor) MyNotifications(ctx context.Context, id uint) ([]*domain.Invitation, error) {
	const op = "uc.notif.my"
	invitations, err := ni.notifRepo.MyNotifications(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return invitations, nil
}
func (ni *NotificationInteractor) AcceptInvite(ctx context.Context, userID, invitation_id, enterprise_id uint) error {
	const op = "uc.notif.accept"
	err := ni.entrRepo.AddUserEnterprise(userID, enterprise_id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if err := ni.notifRepo.DeleteInvitation(ctx, invitation_id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
