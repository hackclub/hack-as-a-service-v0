package main

import (
	"log"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/api"
	"github.com/hackclub/hack-as-a-service/api/oauth"
	"github.com/hackclub/hack-as-a-service/db"
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

	r := gin.Default()

	r.Use(static.Serve("/", static.LocalFile("./frontend/out", false)))

	rg := r.Group("/api")
	err = api.SetupRoutes(rg)
	if err != nil {
		log.Fatalln(err)
	}

	oauth.SetupRoutes(r)

	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/out/404.html")
	})

	r.Run("0.0.0.0:" + getPort())
}
