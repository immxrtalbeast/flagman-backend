package domain

import (
	"context"
	"database/sql/driver"
	"fmt"
	"time"
)

type Role string

const (
	UserRole  Role = "USER"
	AdminRole Role = "ADMIN"
)

// Реализуем интерфейс Scanner для чтения из БД
func (r *Role) Scan(value interface{}) error {
	if value == nil {
		*r = ""
		return nil
	}

	if v, ok := value.([]byte); ok {
		*r = Role(string(v))
		return nil
	}

	return fmt.Errorf("invalid role type: %T", value)
}

// Реализуем интерфейс Valuer для записи в БД
func (r Role) Value() (driver.Value, error) {
	return string(r), nil
}

type User struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	FullName      string `gorm:"unique;not null"`
	Email         string `gorm:"unique;not null"`
	PhoneNumber   string `gorm:"unique;not null"`
	PassHash      []byte `gorm:"not null"`
	CreatedAt     time.Time
	Organizations []Organization `gorm:"many2many:user_organizations;"`
	Roles         []Role         ` gorm:"type:text[];many2many:user_roles;"`

	SentDocuments     []Document          `gorm:"foreignKey:SenderID;references:ID"` // Документы, отправленные пользователем
	ReceivedDocuments []DocumentRecipient `gorm:"foreignKey:UserID"`                 // Документы, полученные для подписи
}

type UserInteractor interface {
	CreateUser(ctx context.Context, fullname string, email string, phonenumber string, pass string) (uint, error)
	Login(ctx context.Context, email string, passhash string) (string, error)
	User(ctx context.Context, id uint) (*User, error)
	// Users(ctx context.Context, page int, limit int) ([]*User, error)
	// UpdateUser(ctx context.Context, id string, name string, login string, passhash string, role Role) error
	// DeleteUser(ctx context.Context, id string) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (uint, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	User(ctx context.Context, id uint) (*User, error)
	// UserByLogin(ctx context.Context, login string) (*User, error)
	// Users(ctx context.Context, page int, limit int) ([]*User, error)
	// UpdateUser(ctx context.Context, user *User) error
	// DeleteUser(ctx context.Context, id string) error
	// UpdateUserPassword(ctx context.Context, passHash []byte) error
}
