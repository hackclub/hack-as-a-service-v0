package apps

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/hackclub/hack-as-a-service/pkg/dokku"
)

func handleGETLogs(c *gin.Context) {
	upgrader := websocket.Upgrader{}

	dokku_conn := c.MustGet("dokkuconn").(*dokku.DokkuConn)
	user := c.MustGet("user").(db.User)

	app_id := c.Param("id")
	tail, err := strconv.Atoi(c.DefaultQuery("tail", "30"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid tail parameter"})
		return
	}

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

	// Open a log stream via the Docker API
	log_stream, err := cli.ContainerLogs(c.Request.Context(), cid, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       strconv.Itoa(tail),
	})
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Spin up a websocket connection
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	defer ws.Close()

	stdout_reader, stdout_writer := io.Pipe()
	stderr_reader, stderr_writer := io.Pipe()

	// Make a channel to hold the log stream
	log_chan := make(chan gin.H)

	// Demultiplex stream in a goroutine
	go func() {
		defer stdout_reader.Close()
		defer stderr_reader.Close()

		defer close(log_chan)

		stdcopy.StdCopy(stdout_writer, stderr_writer, log_stream)
	}()

	// A mutex to ensure stdout and stderr aren't written simultaneously
	log_mutex := sync.Mutex{}

	// Listen for new logs in 2 seperate goroutines for stdout and stderr
	go func() {
		scanner := bufio.NewScanner(stdout_reader)

		for scanner.Scan() {
			log_mutex.Lock()
			log_chan <- gin.H{
				"stream": "stdout",
				"log":    scanner.Text(),
			}
			log_mutex.Unlock()
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr_reader)

		for scanner.Scan() {
			log_mutex.Lock()
			log_chan <- gin.H{
				"stream": "stderr",
				"log":    scanner.Text(),
			}
			log_mutex.Unlock()
		}
	}()

	// Listen for disconnections
	go func() {
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				log_stream.Close()
				return
			}
		}
	}()

	// Close the log stream when done
	defer log_stream.Close()

	for {
		logs, ok := <-log_chan
		if !ok {
			// Channel is closed
			break
		}

		ws.WriteJSON(logs)
	}
}
