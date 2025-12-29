package models

import "time"

type Notification struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null"`
	LoanID    uint   `gorm:"default:null"`                        // Terhubung ke peminjaman (opsional)
	Type      string `gorm:"type:varchar(50);default:'reminder'"` // 'reminder' atau 'denda'
	Message   string `gorm:"type:text"`
	IsRead    bool   `gorm:"default:false"`
	CreatedAt time.Time
}
