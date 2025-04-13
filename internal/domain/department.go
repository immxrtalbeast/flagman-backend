package domain

import (
	"context"
	"time"
)

type Department struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"not null"`
	EnterpriseID uint   `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Permissions  string // JSON с правами (напр., "can_sign", "can_approve")
}

type DepartmentRepository interface {
	// Основные CRUD-операции
	CreateDepartment(ctx context.Context, department *Department) error
	// UpdateDepartment(department *Department) error
	// DeleteDepartment(departmentID uint) error
	// GetDepartmentByID(departmentID uint) (*Department, error)

	// // Управление пользователями в отделе
	// AddUserToDepartment(departmentID uint, userID uint, roleIDs []uint) error
	// RemoveUserFromDepartment(departmentID uint, userID uint) error
	// ListUsersInDepartment(departmentID uint) ([]User, error)

	// // Управление ролями
	// AssignRoleToDepartment(departmentID uint, roleID uint) error
	// RemoveRoleFromDepartment(departmentID uint, roleID uint) error
	// ListRolesInDepartment(departmentID uint) ([]Role, error)
}

type DepartmentInteractor interface {
	CreateDepartment(ctx context.Context, department *Department) error
	// AddUserToDepartment(adminID uint, departmentID uint, userID uint, roles []uint) error
	// RemoveUserFromDepartment(adminID uint, departmentID uint, userID uint) error
	// ListDepartmentUsers(departmentID uint) ([]User, error)

	// AssignRole(adminID uint, departmentID uint, roleID uint) error
	// RevokeRole(adminID uint, departmentID uint, roleID uint) error
}
