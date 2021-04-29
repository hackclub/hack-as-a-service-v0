package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/api/apps"
	"github.com/hackclub/hack-as-a-service/pkg/api/builds"
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

	// uncomment for testing
	// r.GET("/command", func(c *gin.Context) {
	// 	conn := c.MustGet("dokkuconn").(*dokku.DokkuConn)
	// 	cmd := strings.Split(c.Request.URL.Query().Get("command"), " ")
	// 	res, err := conn.RunCommand(c.Request.Context(), cmd)
	// 	if err != nil {
	// 		c.JSON(500, gin.H{"status": "error", "error": err})
	// 	} else {
	// 		c.JSON(200, gin.H{"status": "ok", "output": res})
	// 	}
	// })
	// r.GET("/streamCommand", func(c *gin.Context) {
	// 	conn := c.MustGet("dokkuconn").(*dokku.DokkuConn)
	// 	cmd := strings.Split(c.Request.URL.Query().Get("command"), " ")
	// 	res, err := conn.RunStreamingCommand(c.Request.Context(), cmd)
	// 	if err != nil {
	// 		c.JSON(500, gin.H{"status": "error", "error": err})
	// 	} else {
	// 		c.JSON(200, gin.H{"status": "ok", "output": res})
	// 	}
	// })

	users_rg := r.Group("/users")
	users.SetupRoutes(users_rg)

	teams_rg := r.Group("/teams")
	teams.SetupRoutes(teams_rg)

	apps_rg := r.Group("/apps")
	apps.SetupRoutes(apps_rg)

	builds_rg := r.Group("/builds")
	builds.SetupRoutes(builds_rg)

	return nil
}
