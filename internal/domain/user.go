package domain

import (
	"context"
	"time"
)

type User struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	FullName    string `gorm:"unique;not null"`
	Email       string `gorm:"unique;not null"`
	PhoneNumber string `gorm:"unique;not null"`
	PassHash    []byte `gorm:"not null"`
	CreatedAt   time.Time
	Enterprises []Enterprise `gorm:"many2many:user_enterprises;"`

	SentDocuments     []Document          `gorm:"foreignKey:SenderID;references:ID"` // Документы, отправленные пользователем
	ReceivedDocuments []DocumentRecipient `gorm:"foreignKey:UserID"`                 // Документы, полученные для подписи
}

type UserInteractor interface {
	CreateUser(ctx context.Context, fullname string, email string, phonenumber string, pass string) (uint, error)
	Login(ctx context.Context, email string, passhash string) (string, error)
	User(ctx context.Context, id uint) (*User, error)
	Users(ctx context.Context) ([]*User, error)
	UsersEntr(ctx context.Context, entrID string) ([]User, error)
	// Users(ctx context.Context, enterprisesID []string) ([]*User, error)
	// UpdateUser(ctx context.Context, id string, name string, login string, passhash string, role Role) error
	// DeleteUser(ctx context.Context, id string) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (uint, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	User(ctx context.Context, id uint) (*User, error)
	GetUsersByEnterpriseID(enterpriseID string) ([]User, error)
	Users(ctx context.Context) ([]*User, error)
	UpdateUser(ctx context.Context, user *User) error
	// DeleteUser(ctx context.Context, id string) error
	// UpdateUserPassword(ctx context.Context, passHash []byte) error
}
