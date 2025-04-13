package supabase

import (
	"context"
	"fmt"

	"github.com/immxrtalbeast/flagman-backend/internal/domain"
	"gorm.io/gorm"
)

type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) Create(ctx context.Context, document *domain.Document) error {
	const op = "storage.document.create"
	result := r.db.Create(document)
	if result.Error != nil {
		return fmt.Errorf("%s: %w", op, result.Error)
	}
	return nil
}

func (r *DocumentRepository) DocumentByID(ctx context.Context, id uint) (*domain.Document, error) {
	var document domain.Document
	err := r.db.Where("id = ?", id).Preload("Sender").First(&document).Error
	return &document, err
}
