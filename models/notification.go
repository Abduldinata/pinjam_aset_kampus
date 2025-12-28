package models

import "time"

type Notification struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Message   string    `gorm:"type:text"`
	IsRead    bool      `gorm:"default:false"`
	CreatedAt time.Time
}