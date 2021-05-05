package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/api"
	"github.com/hackclub/hack-as-a-service/pkg/api/auth"
	"github.com/hackclub/hack-as-a-service/pkg/api/oauth"
	"github.com/hackclub/hack-as-a-service/pkg/db"
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

	rg := r.Group("/api", auth.EnsureAuthedUser)
	err = api.SetupRoutes(rg)
	if err != nil {
		log.Fatalln(err)
	}

	gh.SetupRoutes(r.Group("/gh"))

	oauth.SetupRoutes(&r.RouterGroup)

	r.Run("0.0.0.0:" + getPort())
}
