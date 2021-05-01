package db

type User struct {
	Model
	SlackUserID string `gorm:"unique"`
	Name        string
	Avatar      string
	Teams       []*Team `gorm:"many2many:team_users;"`
	Tokens      []Token `json:"-"`
}
