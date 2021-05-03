package main

import (
	"log"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/api"
	"github.com/hackclub/hack-as-a-service/pkg/api/auth"
	"github.com/hackclub/hack-as-a-service/pkg/api/oauth"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/hackclub/hack-as-a-service/pkg/frontend"
	"github.com/hackclub/hack-as-a-service/pkg/gh"
)

func getPort() string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return port
	} else {
		return "5000"
	}
}

func main() {
	err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	if dev := os.Getenv("HAAS_DEV"); dev == "" {
		// prod mode
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Let frontend access cookies in dev
	if dev := os.Getenv("HAAS_DEV"); dev != "" {
		r.Use(func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH")
		})
	}

	r.GET("/swagger.yaml", func(c *gin.Context) { c.File("swagger.yaml") })

	frontend.SetupRoutes(&r.RouterGroup)

	r.Use(static.ServeRoot("/", "./frontend/out"))

	rg := r.Group("/api", auth.EnsureAuthedUser)
	err = api.SetupRoutes(rg)
	if err != nil {
		log.Fatalln(err)
	}

	r.POST("/gh/webhook", gh.HandleWebhook)

	oauth.SetupRoutes(&r.RouterGroup)

	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/out/404.html")
	})

	r.Run("0.0.0.0:" + getPort())
}
