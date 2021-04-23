package db

import "gorm.io/gorm"

type Token struct {
	gorm.Model
	Token  string `gorm:"unique"`
	UserID uint
}
