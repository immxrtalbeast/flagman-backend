package domain

type Organization struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	CreatorID   uint   `gorm:"not null"` // Ссылка на User.ID
	Departments []Department
}
