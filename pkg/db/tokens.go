package db

import (
	"time"
)

type Token struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Token     string `gorm:"primarykey"`
	ExpiresAt time.Time
	UserID    uint
}
