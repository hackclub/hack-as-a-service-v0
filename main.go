package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/dokku"
)

var conn dokku.DokkuConn

func getApiKey() string {
	if key, ok := os.LookupEnv("API_KEY"); ok {
		return key
	} else {
		return "testinghaasapikey"
	}
}

func HandleApi(c *gin.Context) {
	api_key := c.Query("api_key")

	if api_key == "" {
		// Get from auth header if possible
		if auth_header := c.GetHeader("Authorization"); auth_header != "" {
			if strings.HasPrefix(auth_header, "Bearer ") {
				api_key = strings.TrimPrefix(auth_header, "Bearer ")
			}
		}
	}

	if api_key != getApiKey() {
		c.String(401, "Invalid API key")
		return
	}

	args := strings.Split(c.Query("command"), " ")

	res, err := conn.RunCommand(context.Background(), args)
	if err != nil {
		c.String(500, "Error! %s", err)
	} else {
		c.String(200, "Command success:\n%s", res)
	}
}

func getPort() string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return port
	} else {
		return "5000"
	}
}

func main() {
	_conn, err := dokku.DokkuConnect(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	conn = _conn

	r := gin.Default()

	r.Use(static.Serve("/", static.LocalFile("./frontend/out", false)))
	r.GET("/api", HandleApi)
	rg := r.Group("/api")
	err = dokku.SetupRoutes(rg)
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
