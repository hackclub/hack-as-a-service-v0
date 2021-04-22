package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	SlackUserID string
}

type BillingAccount struct {
	gorm.Model
	HNUserID string
}

type App struct {
	gorm.Model
	Name             string
	BillingAccountID int
	BillingAccount   BillingAccount
	Users            []User `gorm:"many2many:user_apps;"`
	// TODO: add more fields
}
