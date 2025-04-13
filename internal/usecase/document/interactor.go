package document

import (
	"context"
	"fmt"
	"time"

	"github.com/immxrtalbeast/flagman-backend/internal/domain"
)

type DocumentInteractor struct {
	documentRepo domain.DocumentRepository
	usrRepo      domain.UserRepository
}

func NewDocumentInteractor(documentRepo domain.DocumentRepository, usrRepo domain.UserRepository) *DocumentInteractor {
	return &DocumentInteractor{
		documentRepo: documentRepo,
		usrRepo:      usrRepo,
	}
}
func (di *DocumentInteractor) CreateDocument(ctx context.Context, senderID uint, title, filePath string) (*domain.Document, error) {
	const op = "uc.document.create"
	sender, err := di.usrRepo.User(ctx, senderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	document := domain.Document{
		Title:     title,
		Sender:    *sender,
		FilePath:  filePath,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := di.documentRepo.Create(ctx, &document); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &document, nil
}
func (di *DocumentInteractor) DocumentByID(ctx context.Context, id uint) (*domain.Document, error) {
	const op = "uc.document.byID"
	document, err := di.documentRepo.DocumentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return document, nil
}
