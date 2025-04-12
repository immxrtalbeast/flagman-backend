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

// Получатель документа (связь между документом и пользователем)
type DocumentRecipient struct {
	ID         uint       `gorm:"primaryKey;autoIncrement"`
	DocumentID uint       // Ссылка на документ
	UserID     uint       // ID получателя
	Status     string     `gorm:"default:'pending'"` // Статус для получателя: "pending", "signed", "rejected"
	Signature  string     // Уникальная подпись (хеш телефона + соль)
	SignedAt   *time.Time // Время подписания
}

type DocumentRepository interface {
	Create(ctx context.Context, document *Document) error
	// FindByID(id uuid.UUID) (*Document, error)
	// Update(document *Document) error
	// Delete(id uuid.UUID) error
	// FindBySender(senderID uint) ([]Document, error)
}

type DocumentInteractor interface {
	CreateDocument(ctx context.Context, senderID uint, title, filePath string) (*Document, error)
	// SendDocument(documentID uuid.UUID, recipientIDs []uint) error
	// GetDocument(documentID uuid.UUID) (*Document, error)
	// UpdateDocument(document *Document) error
	// DeleteDocument(documentID uuid.UUID) error
	// ListUserDocuments(senderID uint) ([]Document, error)
	// GetDocumentRecipients(documentID uuid.UUID) ([]DocumentRecipient, error)
}

// type DocumentRecipientRepository interface {
// 	Create(recipient *DocumentRecipient) error
// 	FindByID(id uuid.UUID) (*DocumentRecipient, error)
// 	Update(recipient *DocumentRecipient) error
// 	Delete(id uuid.UUID) error
// 	FindByDocument(documentID uuid.UUID) ([]DocumentRecipient, error)
// 	FindByUser(userID uint) ([]DocumentRecipient, error)
// }

// type DocumentRecipientUsecase interface {
// 	SignDocument(recipientID uuid.UUID, signature string) error
// 	RejectDocument(recipientID uuid.UUID) error
// 	GetRecipient(recipientID uuid.UUID) (*DocumentRecipient, error)
// 	ListUserDocuments(userID uint) ([]DocumentRecipient, error)
// 	GetDocumentStatus(documentID uuid.UUID) (string, error)
// }
