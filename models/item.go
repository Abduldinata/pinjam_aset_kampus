package models

import (
	"time"
)

type Item struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"type:varchar(100)"`
	Category    string    `gorm:"type:varchar(50)"`
	Stock       int       `gorm:"not null;default:0"`
	Location    string    `gorm:"type:varchar(100)"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}