package db

import (
	"time"

	"gorm.io/gorm"
)

type Token struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Token     string         `gorm:"primarykey"`
	ExpiresAt time.Time
	UserID    uint
}
