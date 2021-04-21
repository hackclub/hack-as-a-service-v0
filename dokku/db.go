package dokku

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
	name           string
	billingAccount BillingAccount
	users          []User
	// TODO: add more fields
}
