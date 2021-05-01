package apps

import (
	"bufio"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/api/util"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/hackclub/hack-as-a-service/pkg/dokku"
)

type Stats struct {
	Read     time.Time `json:"read"`
	CpuStats struct {
		CpuUsage struct {
			PercpuUsage []interface{} `json:"percpu_usage"`
			TotalUsage  int64         `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemCpuUsage int64
		OnlineCpus     int64
	} `json:"cpu_stats"`
	PrecpuStats struct {
		CpuUsage struct {
			TotalUsage int64 `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemCpuUsage int64 `json:"system_cpu_usage"`
		OnlineCpus     int64 `json:"online_cpus"`
	} `json:"precpu_stats"`
	MemoryStats struct {
		Stats struct {
			Cache int64 `json:"cache"`
		} `json:"stats"`
		Usage int64 `json:"usage"`
		Limit int64 `json:"limit"`
	} `json:"memory_stats"`
}

type WsOutput struct {
	Timestamp       time.Time
	MemUsagePercent float64
	CpuUsagePercent float64
}

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
		var stat Stats
		if err := json.Unmarshal([]byte(line), &stat); err != nil {
			log.Printf("Error decoding json: %+v\n", err)
			break
		}
		// log.Printf("Stat = %+v\n", stat)
		// https://docs.docker.com/engine/api/v1.41/#operation/ContainerStats
		mem_usage := (float64(stat.MemoryStats.Usage) - float64(stat.MemoryStats.Stats.Cache)) / float64(stat.MemoryStats.Limit) * 100.
		cpu_usage := float64(stat.CpuStats.CpuUsage.TotalUsage) - float64(stat.PrecpuStats.CpuUsage.TotalUsage)
		x := float64(stat.CpuStats.SystemCpuUsage) - float64(stat.PrecpuStats.SystemCpuUsage)
		if x != 0 {
			cpu_usage /= x
		}
		cpu_usage *= float64(len(stat.CpuStats.CpuUsage.PercpuUsage)) * 100.
		output := WsOutput{
			Timestamp:       stat.Read,
			MemUsagePercent: mem_usage,
			CpuUsagePercent: cpu_usage,
		}
		// log.Printf("Output = %+v\n", output)
		err := ws.WriteJSON(output)
		if err != nil {
			log.Printf("Error writing to ws: %+v\n", err)
			break
		}
	}
}
