package domain

import "time"

type Enterprise struct {
	ID          uint      `gorm:"primaryKey"`      // Уникальный идентификатор
	Name        string    `gorm:"unique;not null"` // Название организации (уникальное)
	Description string    // Описание (опционально)
	CreatorID   uint      `gorm:"not null"` // ID создателя (User.ID)
	CreatedAt   time.Time // Дата создания
	UpdatedAt   time.Time // Дата обновления
	Users       []User    `gorm:"many2many:user_enterprises;"`
	// Связи
	Departments []Department `gorm:"foreignKey:EnterpriseID"` // Отделы в организации
}

type EnterpriseRepository interface {
	// Основные CRUD-операции
	CreateEnterprise(enterprise *Enterprise) error
	AddUserEnterprise(userID uint, enterpriseID uint) error
	EnterpriseByID(enterpriseID uint) (*Enterprise, error)
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
