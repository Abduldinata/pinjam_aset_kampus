package models

import (
	"time"
)

type Loan struct {
	ID            uint       `gorm:"primaryKey"`
	UserID        uint       `gorm:"not null"`
	ItemID        uint       `gorm:"not null"`
	BorrowDate    time.Time  `gorm:"type:date"`
	DueDate       time.Time  `gorm:"type:date"`
	ReturnDate    *time.Time `gorm:"type:date"` // Pakai pointer (*) karena bisa NULL kalau belum kembali
	Status        string     `gorm:"type:status_type;default:'dipinjam'"`
	Notes         string     `gorm:"type:text"`
	FineAmount    int        `gorm:"default:0"`        // Total denda yang harus dibayar
	IsFinePaid    bool       `gorm:"default:false"`    // Status kelunasan denda
	PaymentMethod string     `gorm:"type:varchar(50)"` // DANA, Bank, dll
	PaymentProof  string     `gorm:"type:text"`        // Path/nama file bukti transfer
	StudentIDCard string     `gorm:"type:text"`        // Path/nama file KTM saat pinjam
	CreatedAt     time.Time
	UpdatedAt     time.Time

	// Relasi (Associations) - Agar bisa ambil data User dan Item terkait
	User User `gorm:"foreignKey:UserID"`
	Item Item `gorm:"foreignKey:ItemID"`
}
