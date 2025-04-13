package supabase

import (
	"context"
	"fmt"
	"time"

	"github.com/immxrtalbeast/flagman-backend/internal/domain"
	"gorm.io/gorm"
)

type DocumentRecipientRepository struct {
	db *gorm.DB
}

func NewDocumentRecipientRepository(db *gorm.DB) *DocumentRecipientRepository {
	return &DocumentRecipientRepository{db: db}
}

func (r *DocumentRecipientRepository) Create(recipient *domain.DocumentRecipient) error {
	const op = "storage.documentrecipi.create"
	result := r.db.Create(recipient)
	if result.Error != nil {
		return fmt.Errorf("%s: %w", op, result.Error)
	}
	return nil
}
func (r *DocumentRecipientRepository) FindByID(id string) (*domain.DocumentRecipient, error) {
	var recipient domain.DocumentRecipient
	err := r.db.Where("id = ?", id).Preload("Document.Sender").First(&recipient).Error

	return &recipient, err
}

func (r *DocumentRecipientRepository) SignDocument(ctx context.Context, id string, userID uint, signature string) error {
	var recipient domain.DocumentRecipient
	r.db.Where("id = ? AND user_id = ?", id, userID).First(&recipient)
	if recipient.Status != "pending" {
		return fmt.Errorf("Already signed or declined.")
	}
	now := time.Now()
	err := r.db.Model(&recipient).Updates(map[string]interface{}{
		"status":     "signed",
		"signature":  signature,
		"signed_at":  now,
		"updated_at": now,
	})
	return err.Error
}
func (r *DocumentRecipientRepository) RejectDocument(ctx context.Context, id string, userID uint) error {
	var recipient domain.DocumentRecipient
	r.db.Where("id = ? AND user_id = ?", id, userID).First(&recipient)
	if recipient.Status != "pending" {
		fmt.Errorf("Already signed or declined.")
	}
	now := time.Now()
	err := r.db.Model(&recipient).Updates(map[string]interface{}{
		"status":     "rejected",
		"signed_at":  now,
		"updated_at": now,
	})
	return err.Error
}

func (r *DocumentRecipientRepository) ListUserDocuments(ctx context.Context, userID uint, status string) ([]domain.DocumentRecipient, error) {
	var documents []domain.DocumentRecipient
	if status == "" {
		err := r.db.Where("user_id = ?", userID).Preload("Document.Sender").WithContext(ctx).Find(&documents).Error
		return documents, err
	} else {
		err := r.db.Where("user_id = ? AND status = ?", userID, status).Preload("Document.Sender").WithContext(ctx).Find(&documents).Error
		return documents, err
	}

}
