package domain

import (
	"context"
	"time"
)

type Invitation struct {
	ID             uint   `gorm:"primaryKey"`
	Email          string `gorm:"not null"` // Email приглашенного
	ReceiverID     uint
	EnterpriseName string
	EnterpriseID   uint // ID организации
	CreatedBy      uint // ID приглашающего
}

type Enterprise struct {
	ID          uint      `gorm:"primaryKey"`      // Уникальный идентификатор
	Name        string    `gorm:"unique;not null"` // Название организации (уникальное)
	Description string    // Описание (опционально)
	CreatorID   uint      `gorm:"not null"` // ID создателя (User.ID)
	CreatedAt   time.Time // Дата создания
	UpdatedAt   time.Time // Дата обновления
	Users       []User    `gorm:"many2many:user_enterprises;"`
	// Связи
	// Departments []Department `gorm:"foreignKey:EnterpriseID"` // Отделы в организации
}

type EnterpriseRepository interface {
	// Основные CRUD-операции
	CreateEnterprise(enterprise *Enterprise) error
	AddUserEnterprise(userID uint, enterpriseID uint) error
	EnterpriseByID(enterpriseID uint) (*Enterprise, error)
	GetEnterprisesByUserID(userID uint) ([]Enterprise, error)
	// UpdateEnterprise(enterprise *Enterprise) error
	// DeleteEnterprise(enterpriseID uint) error
	// GetEnterpriseByID(enterpriseID uint) (*Enterprise, error)
	// ListEnterprisesByUser(userID uint) ([]Enterprise, error) // Организации, где состоит пользователь

	// // Управление отделами
	// AddDepartmentToEnterprise(enterpriseID uint, department *Department) error
	// RemoveDepartmentFromEnterprise(enterpriseID uint, departmentID uint) error
	// ListDepartments(enterpriseID uint) ([]Department, error)

	// // Приглашение пользователей
	// InviteUserToEnterprise(enterpriseID uint, userID uint) error
	// RemoveUserFromEnterprise(enterpriseID uint, userID uint) error
}

type EnterpriseInteractor interface {
	CreateEnterprise(userID uint, name string, description string) (*Enterprise, error)
	EnterpriseByID(enterpriseID uint) (*Enterprise, error)
	InviteUser(senderID uint, userEmail string, enterpriseID uint, enterspriseName string) (*Invitation, error)
	GetEnterprisesByUserID(userID uint) ([]Enterprise, error)
	// UpdateEnterprise(userID uint, enterprise *Enterprise) error // Проверка прав
	// DeleteEnterprise(userID uint, enterpriseID uint) error
	// GetEnterprise(userID uint, enterpriseID uint) (*Enterprise, error) // Проверка доступа

	// // Отделы
	// CreateDepartment(userID uint, enterpriseID uint, name string) (*Department, error)
	// DeleteDepartment(userID uint, departmentID uint) error

	// // Пользователи
	// InviteUserToEnterprise(inviterID uint, enterpriseID uint, email string) error
	// RemoveUserFromEnterprise(adminID uint, enterpriseID uint, userID uint) error
}
type NotificationRepository interface {
	CreateInvitation(ctx context.Context, invite Invitation) error
	MyNotifications(ctx context.Context, id uint) ([]*Invitation, error)
	DeleteInvitation(ctx context.Context, id uint) error
}
type NotificationInteractor interface {
	MyNotifications(ctx context.Context, id uint) ([]*Invitation, error)
	AcceptInvite(ctx context.Context, userID, invitation_id, enterprise_id uint) error
}
