package db

import (
	"github.com/shopspring/decimal"
)

type Team struct {
	Model
	Name      string
	Avatar    string
	Automatic bool            // Whether the team was created automatically for ad-hoc app sharing
	Personal  bool            // Whether this is a user's personal team
	Expenses  decimal.Decimal `gorm:"type:decimal(64, 18);default:0"` // Outstanding expenses for this team
	Users     []*User         `gorm:"many2many:team_users;"`
	Apps      []App
}
