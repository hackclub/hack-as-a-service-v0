package db

import "gorm.io/gorm"

type Team struct {
	gorm.Model
	Name      string
	Automatic bool
	HNUserID  string
	Users     []User `gorm:"many2many:team_users;"`
}
