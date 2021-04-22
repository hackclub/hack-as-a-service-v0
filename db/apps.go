package db

import "gorm.io/gorm"

type App struct {
	gorm.Model
	Name             string
	BillingAccountID int
	BillingAccount   BillingAccount
	Users            []User `gorm:"many2many:user_apps;"`
	// TODO: add more fields
}
