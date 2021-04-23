package db

import (
	"time"

	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	Token     string `gorm:"unique"`
	ExpiresAt time.Time
	UserID    uint
}
