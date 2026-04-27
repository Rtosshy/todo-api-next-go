package dao

type StatusID int

type StatusName string

type Status struct {
	ID   StatusID   `gorm:"primaryKey"`
	Name StatusName `gorm:"not null"`
}
