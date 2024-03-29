package db

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Model is a clone of a gorm.Model, but without the `DeletedAt` field
type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var DB *gorm.DB

func getGsn() string {
	if db_host := os.Getenv("DATABASE_URL"); db_host != "" {
		return db_host
	}
	db_password := os.Getenv("POSTGRES_PASSWORD")
	return fmt.Sprintf("host=db user=postgres password=%s dbname=haas port=5432 sslmode=disable", db_password)
}

func Connect() error {
	_db, err := gorm.Open(postgres.Open(getGsn()), &gorm.Config{})
	if err != nil {
		return err
	}

	err = _db.AutoMigrate(&User{}, &Team{}, &App{}, &Token{}, &Build{})

	if err != nil {
		return err
	}

	DB = _db

	return nil
}
