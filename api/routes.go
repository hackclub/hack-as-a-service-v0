package api

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/api/apps"
	"github.com/hackclub/hack-as-a-service/api/billing"
	"github.com/hackclub/hack-as-a-service/api/users"
	"github.com/hackclub/hack-as-a-service/db"
	"github.com/hackclub/hack-as-a-service/dokku"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getGsn() string {
	db_password := os.Getenv("POSTGRES_PASSWORD")
	db_host := "db"
	if db_host2, ok := os.LookupEnv("DATABASE_URL"); ok {
		db_host = db_host2
	}
	return fmt.Sprintf("host=%s user=postgres password=%s dbname=haas port=5432 sslmode=disable", db_host, db_password)
}

func SetupRoutes(r *gin.RouterGroup) error {
	conn, err := dokku.DokkuConnect(context.Background())
	if err != nil {
		return err
	}
	r.Use(func(c *gin.Context) {
		c.Set("dokkuconn", conn)
	})

	r.GET("/", handleApiCommand)

	_db, err := gorm.Open(postgres.Open(getGsn()), &gorm.Config{})
	if err != nil {
		return err
	}
	err = _db.AutoMigrate(&db.User{})
	if err != nil {
		return err
	}
	err = _db.AutoMigrate(&db.BillingAccount{})
	if err != nil {
		return err
	}
	err = _db.AutoMigrate(&db.App{})
	if err != nil {
		return err
	}
	r.Use(func(c *gin.Context) {
		c.Set("db", _db)
	})

	users_rg := r.Group("/users")
	users.SetupRoutes(users_rg)
	billingAccounts_rg := r.Group("/billingAccounts")
	billing.SetupRoutes(billingAccounts_rg)
	apps_rg := r.Group("/apps")
	apps.SetupRoutes(apps_rg)

	return nil
}
