package models

import (
	"time"
)

type Loan struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null"`
	ItemID     uint      `gorm:"not null"`
	BorrowDate time.Time `gorm:"type:date"`
	DueDate    time.Time `gorm:"type:date"`
	ReturnDate *time.Time `gorm:"type:date"` // Pakai pointer (*) karena bisa NULL kalau belum kembali
	Status     string    `gorm:"type:status_type;default:'dipinjam'"`
	Notes      string    `gorm:"type:text"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Relasi (Associations) - Agar bisa ambil data User dan Item terkait
	User User `gorm:"foreignKey:UserID"`
	Item Item `gorm:"foreignKey:ItemID"`
}