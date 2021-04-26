package db

import "gorm.io/gorm"

type Team struct {
	gorm.Model
	Name      string
	Automatic bool // Whether the team was created automatically for ad-hoc app sharing
	Personal  bool // Whether this is a user's personal team
	HNUserID  string
	Users     []*User `gorm:"many2many:team_users;"`
}