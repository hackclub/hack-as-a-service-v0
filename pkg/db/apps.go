package db

type App struct {
	Model
	Name string
	// The app's Dokku name
	ShortName string
	TeamID    uint
	// TODO: add more fields
}
