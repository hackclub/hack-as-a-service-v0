package db

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func get_gsn() string {
	db_password := os.Getenv("POSTGRES_PASSWORD")
	return fmt.Sprintf("host=db user=postgres password=%s dbname=haas port=5432 sslmode=disable", db_password)
}

func SetupRoutes(r *gin.RouterGroup) error {
	_db, err := gorm.Open(postgres.Open(get_gsn()), &gorm.Config{})
	if err != nil {
		return err
	}
	db = _db
	err = db.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	users_rg := r.Group("/users")
	setupUserRoutes(users_rg)

	return nil
}
