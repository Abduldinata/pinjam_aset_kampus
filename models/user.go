package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(100)"`
	Email     string    `gorm:"unique;type:varchar(100)"`
	Password  string    `gorm:"type:varchar(255)"`
	Role      string    `gorm:"type:role_type;default:'user'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}