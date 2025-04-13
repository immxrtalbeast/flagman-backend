package domain

import (
	"context"
	"time"
)

// Документ
type Document struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Title     string // Название документа
	FilePath  string // Путь к файлу в S3 (например: "documents/uuid.pdf")
	Status    string // Общий статус: "draft", "sent", "archived"
	SenderID  uint
	Sender    User `gorm:"foreignKey:SenderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Связи
	Recipients []DocumentRecipient `gorm:"foreignKey:DocumentID"` // Получатели
}

type DocumentRecipient struct {
	ID         uint     `gorm:"primaryKey;autoIncrement"`
	DocumentID uint     `gorm:"index"`
	Document   Document `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID     uint     `gorm:"index"`
	Status     string   `gorm:"default:'pending'"`
	Signature  string
	SignedAt   *time.Time
}
type DocumentRepository interface {
	Create(ctx context.Context, document *Document) error
	DocumentByID(ctx context.Context, id uint) (*Document, error)
	// Update(document *Document) error
	// Delete(id uuid.UUID) error
	// FindBySender(senderID uint) ([]Document, error)
}

type DocumentInteractor interface {
	CreateDocument(ctx context.Context, senderID uint, title, filePath string) (*Document, error)
	DocumentByID(ctx context.Context, id uint) (*Document, error)
	// SendDocument(documentID uuid.UUID, recipientIDs []uint) error
	// GetDocument(documentID uuid.UUID) (*Document, error)
	// UpdateDocument(document *Document) error
	// DeleteDocument(documentID uuid.UUID) error
	// ListUserDocuments(senderID uint) ([]Document, error)
	// GetDocumentRecipients(documentID uuid.UUID) ([]DocumentRecipient, error)
}

type DocumentRecipientRepository interface {
	Create(recipient *DocumentRecipient) error
	FindByID(id string) (*DocumentRecipient, error)
	SignDocument(ctx context.Context, id string, userID uint, signature string) error
	RejectDocument(ctx context.Context, id string, userID uint) error
	ListUserDocuments(ctx context.Context, userID uint, status string) ([]DocumentRecipient, error)
	// Update(recipient *DocumentRecipient) error
	// Delete(id uuid.UUID) error
	// FindByDocument(documentID uuid.UUID) ([]DocumentRecipient, error)
	// FindByUser(userID uint) ([]DocumentRecipient, error)
}

type DocumentRecipientInteractor interface {
	SignDocument(ctx context.Context, recipientID string, userID uint) error
	RejectDocument(ctx context.Context, id string, userID uint) error
	// GetRecipient(recipientID uuid.UUID) (*DocumentRecipient, error)
	ListUserDocuments(ctx context.Context, userID uint, status string) ([]DocumentRecipient, error)
	// GetDocumentStatus(documentID uuid.UUID) (string, error)
}
