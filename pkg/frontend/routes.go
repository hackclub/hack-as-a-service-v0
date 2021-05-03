package frontend

import (
	"path"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func sendHTML(file string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.File(path.Join("frontend", "out", file))
	}
}

func SetupRoutes(r *gin.RouterGroup) {
	r.GET("/", sendHTML("index.html"))
	r.GET("/swagger", sendHTML("swagger.html"))
	r.GET("/landing", sendHTML("landing.html"))
	r.GET("/dashboard", sendHTML("dashboard.html"))
	r.GET("/teams/:id", sendHTML("teams/[id].html"))
	r.GET("/apps/:id", sendHTML("apps/[id].html"))
	r.GET("/apps/:id/logs", sendHTML("apps/[id]/logs.html"))
	r.Use(static.ServeRoot("../../assets", "./assets"))
}
