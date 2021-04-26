package db

import "gorm.io/gorm"

type App struct {
	gorm.Model
	Name string
	// The app's Dokku name
	ShortName string
	TeamID    uint
	Team      Team
	// TODO: add more fields
}
