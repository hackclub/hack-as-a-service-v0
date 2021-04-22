package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/api"
)

func getPort() string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return port
	} else {
		return "5000"
	}
}

func main() {
	r := gin.Default()

	r.Use(static.Serve("/", static.LocalFile("./frontend/out", false)))
	rg := r.Group("/api", api.RequireBearerAuth())
	err := api.SetupRoutes(rg)
	if err != nil {
		log.Fatalln(err)
	}

	r.GET("/oauth/login", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("https://slack.com/oauth/v2/authorize?user_scope=identity.basic&client_id=%s", os.Getenv("SLACK_CLIENT_ID")))
	})

	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/out/404.html")
	})

	r.Run("0.0.0.0:" + getPort())
}
