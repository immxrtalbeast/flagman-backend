package domain

type Department struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"not null"`
	OrganizationID uint   `gorm:"not null"`
	Permissions    string // JSON с правами (напр., "can_sign", "can_approve")
}
