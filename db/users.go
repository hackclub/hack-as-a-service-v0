package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	SlackUserID string `gorm:"unique"`
	Name        string
	Avatar      string
	Tokens      []Token
}
