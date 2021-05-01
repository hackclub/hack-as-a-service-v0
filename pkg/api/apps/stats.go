package apps

import (
	"bufio"
	"encoding/json"
	"log"
	"strings"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/api/util"
	"github.com/hackclub/hack-as-a-service/pkg/biller"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/hackclub/hack-as-a-service/pkg/dokku"
)

func handleGETStats(c *gin.Context) {
	upgrader := util.MakeWebsocketUpgrader()

	dokku_conn := c.MustGet("dokkuconn").(*dokku.DokkuConn)
	user := c.MustGet("user").(db.User)

	app_id := c.Param("id")

	var app db.App

	result := db.DB.Joins("JOIN teams ON teams.id = apps.team_id").
		Joins("JOIN team_users ON team_users.team_id = teams.id").
		First(&app, "apps.id = ? AND team_users.user_id = ?", app_id, user.ID)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	// Get the app's container ID
	cid, err := dokku_conn.RunCommand(c.Request.Context(), []string{"haas:cid", app.ShortName})
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	cid = strings.TrimSpace(cid)

	// Initialize a Docker API client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": "Error creating Docker client"})
		return
	}

	stats_stream, err := cli.ContainerStats(c.Request.Context(), cid, true)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": "Error getting container logs"})
		return
	}
	lines := bufio.NewScanner(stats_stream.Body)

	// Spin up a websocket connection
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	defer stats_stream.Body.Close()
	defer ws.Close()

	// Drop the first line since it contains bad data
	// FIXME: why?
	if !lines.Scan() {
		return
	}
	for lines.Scan() {
		line := lines.Text()
		// log.Printf("Got line: %s\n", line)
		var stat biller.Stats
		if err := json.Unmarshal([]byte(line), &stat); err != nil {
			log.Printf("Error decoding json: %+v\n", err)
			break
		}
		output := stat.Process()
		// log.Printf("Output = %+v\n", output)
		err := ws.WriteJSON(output)
		if err != nil {
			log.Printf("Error writing to ws: %+v\n", err)
			break
		}
	}
}
