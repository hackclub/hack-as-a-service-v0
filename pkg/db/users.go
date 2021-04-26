package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	SlackUserID string `gorm:"unique"`
	Name        string
	Avatar      string
	Teams       []*Team `gorm:"many2many:team_users;"`
	Tokens      []Token `json:"-"`
}
