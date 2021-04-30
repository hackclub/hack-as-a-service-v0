package frontend

import (
	"path"

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
	r.GET("/team/:id", sendHTML("team/[id].html"))
	r.GET("/app/:id", sendHTML("app/[id].html"))
}
