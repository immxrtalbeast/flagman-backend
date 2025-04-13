package domain

// import (
// 	"context"
// 	"time"
// )

// type Role struct {
// 	ID          uint   `gorm:"primaryKey"`
// 	Name        string `gorm:"unique;not null"` // Например: "accountant", "manager"
// 	Permissions string `gorm:"type:json"`       // JSON: {"allowed_departments": [1, 2, 3]}
// }

// type UserDepartment struct {
// 	UserID       uint `gorm:"primaryKey"`
// 	DepartmentID uint `gorm:"primaryKey"`
// 	CreatedAt    time.Time
// 	UpdatedAt    time.Time

// 	User       User       `gorm:"foreignKey:UserID"`
// 	Department Department `gorm:"foreignKey:DepartmentID"`
// }

// type Department struct {
// 	ID           uint   `gorm:"primaryKey"`
// 	Name         string `gorm:"not null"`
// 	RoleID       uint   `gorm:"not null"`
// 	EnterpriseID uint   `gorm:"not null"`
// 	Users        User   `gorm:"many2many:user_departments;"`
// 	Role         Role   `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
// 	CreatedAt    time.Time
// 	UpdatedAt    time.Time
// }

// // RoleRepository - работа с ролями
// type RoleRepository interface {
// 	CreateRole(ctx context.Context, role *Role) error
// 	// UpdateRole(ctx context.Context, role *Role) error
// 	// DeleteRole(ctx context.Context, roleID uint) error
// 	// GetRoleByID(ctx context.Context, roleID uint) (*Role, error)
// 	// ListRoles(ctx context.Context) ([]Role, error)
// }

// // DepartmentRepository - работа с отделами
// type DepartmentRepository interface {
// 	CreateDepartment(ctx context.Context, department *Department) error
// 	UpdateDepartment(ctx context.Context, department *Department) error
// 	DeleteDepartment(ctx context.Context, departmentID uint) error
// 	GetDepartmentByID(ctx context.Context, departmentID uint) (*Department, error)
// 	ListDepartmentsByEnterprise(ctx context.Context, enterpriseID uint) ([]Department, error)
// }

// // RoleInteractor - бизнес-логика для ролей
// type RoleInteractor interface {
// 	CreateRole(ctx context.Context, name, permissions string) (*Role, error)
// 	// UpdateRole(ctx context.Context, roleID uint, name, permissions string) error
// 	// DeleteRole(ctx context.Context, roleID uint) error
// 	// GetRole(ctx context.Context, roleID uint) (*Role, error)
// }

// // DepartmentInteractor - бизнес-логика для отделов
// type DepartmentInteractor interface {
// 	CreateDepartment(
// 		ctx context.Context,
// 		name string,
// 		roleID uint,
// 		enterpriseID uint,
// 	) (*Department, error)
// 	UpdateDepartmentRole(ctx context.Context, departmentID uint, newRoleID uint) error
// 	GetDepartment(ctx context.Context, departmentID uint) (*Department, error)
// }

// // PermissionInteractor - проверка разрешений
// type PermissionInteractor interface {
// 	CanPerformAction(
// 		ctx context.Context,
// 		roleID uint,
// 		action string,
// 		targetDepartmentID uint,
// 	) (bool, error)
// }
