package db

import "gorm.io/gorm"

type App struct {
	gorm.Model
	Name   string
	TeamID uint
	Team   Team
	// TODO: add more fields
}
