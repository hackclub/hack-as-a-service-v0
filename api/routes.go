package api

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/api/apps"
	"github.com/hackclub/hack-as-a-service/api/billing"
	"github.com/hackclub/hack-as-a-service/api/users"
	"github.com/hackclub/hack-as-a-service/dokku"
)

func getGsn() string {
	if db_host := os.Getenv("DATABASE_URL"); db_host != "" {
		return db_host
	}
	db_password := os.Getenv("POSTGRES_PASSWORD")
	return fmt.Sprintf("host=db user=postgres password=%s dbname=haas port=5432 sslmode=disable", db_password)
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

	users_rg := r.Group("/users")
	users.SetupRoutes(users_rg)

	billingAccounts_rg := r.Group("/billingAccounts")
	billing.SetupRoutes(billingAccounts_rg)

	apps_rg := r.Group("/apps")
	apps.SetupRoutes(apps_rg)

	return nil
}
