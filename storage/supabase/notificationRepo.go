package supabase

import (
	"context"

	"github.com/immxrtalbeast/flagman-backend/internal/domain"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) CreateInvitation(ctx context.Context, invite domain.Invitation) error {
	result := r.db.Create(&invite)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *NotificationRepository) DeleteInvitation(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Invitation{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *NotificationRepository) MyNotifications(ctx context.Context, id uint) ([]*domain.Invitation, error) {
	var invitations []*domain.Invitation
	err := r.db.WithContext(ctx).Where("receiver_id = ?", id).Find(&invitations).Error
	return invitations, err
}
