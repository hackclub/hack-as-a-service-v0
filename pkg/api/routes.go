package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/api/apps"
	"github.com/hackclub/hack-as-a-service/pkg/api/teams"
	"github.com/hackclub/hack-as-a-service/pkg/api/users"
	"github.com/hackclub/hack-as-a-service/pkg/dokku"
)

func SetupRoutes(r *gin.RouterGroup) error {
	conn, err := dokku.DokkuConnect(context.Background())
	if err != nil {
		return err
	}
	r.Use(func(c *gin.Context) {
		c.Set("dokkuconn", conn)
	})

	users_rg := r.Group("/users")
	users.SetupRoutes(users_rg)

	teams_rg := r.Group("/teams")
	teams.SetupRoutes(teams_rg)

	apps_rg := r.Group("/apps")
	apps.SetupRoutes(apps_rg)

	return nil
}
