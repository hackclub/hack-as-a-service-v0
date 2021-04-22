package db

import "gorm.io/gorm"

type BillingAccount struct {
	gorm.Model
	HNUserID string
}
