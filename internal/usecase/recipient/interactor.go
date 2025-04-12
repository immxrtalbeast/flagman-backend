package recipient

import (
	"context"
	"fmt"

	"github.com/immxrtalbeast/flagman-backend/internal/domain"
	"github.com/immxrtalbeast/flagman-backend/internal/lib"
)

type DocumentRecipientInteractor struct {
	recipiRepo domain.DocumentRecipientRepository
	usrRepo    domain.UserRepository
	secreSalt  string
}

func NewDocumentRecipientInteractor(recipiRepo domain.DocumentRecipientRepository, usrRepo domain.UserRepository, secretSalt string) *DocumentRecipientInteractor {
	return &DocumentRecipientInteractor{
		recipiRepo: recipiRepo,
		usrRepo:    usrRepo,
		secreSalt:  secretSalt,
	}
}

func (dr *DocumentRecipientInteractor) SignDocument(ctx context.Context, recipientID string, userID uint) error {
	const op = "uc.recipient.sign"
	user, err := dr.usrRepo.User(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	signature := lib.NewSignature(user, dr.secreSalt)

	if err := dr.recipiRepo.SignDocument(ctx, recipientID, user.ID, signature); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
func (dr *DocumentRecipientInteractor) RejectDocument(ctx context.Context, recipientID string, userID uint) error {
	const op = "uc.recipient.reject"
	err := dr.recipiRepo.RejectDocument(ctx, recipientID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (dr *DocumentRecipientInteractor) ListUserDocuments(ctx context.Context, userID uint, status string) ([]domain.DocumentRecipient, error) {
	const op = "uc.recipient.list"

	list, err := dr.recipiRepo.ListUserDocuments(ctx, userID, status)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return list, nil
}
