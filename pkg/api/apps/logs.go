package apps

import (
	"bufio"
	"io"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/hackclub/hack-as-a-service/pkg/dokku"
)

func handleGETLogs(c *gin.Context) {
	dokku_conn := c.MustGet("dokkuconn").(*dokku.DokkuConn)
	user := c.MustGet("user").(db.User)

	app_id := c.Param("id")

	var app db.App
	result := db.DB.Raw(`SELECT apps.* FROM apps
	JOIN teams ON teams.id = apps.team_id
	JOIN team_users ON team_users.team_id = teams.id
	WHERE apps.id = ? AND team_users.user_id = ?`, app_id, user.ID).Scan(&app)
	if result.Error != nil {
		c.JSON(404, gin.H{"status": "error", "message": result.Error.Error()})
		return
	} else if result.RowsAffected < 1 {
		c.JSON(404, gin.H{"status": "error", "message": "App not found"})
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
		Tail:       "30",
	})
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	stdout_reader, stdout_writer := io.Pipe()
	stderr_reader, stderr_writer := io.Pipe()

	// Make a channel to hold the log stream
	log_chan := make(chan []byte)

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
			log_chan <- append([]byte("[stdout] "), scanner.Bytes()...)
			log_mutex.Unlock()
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr_reader)

		for scanner.Scan() {
			log_mutex.Lock()
			log_chan <- append([]byte("[stderr] "), scanner.Bytes()...)
			log_mutex.Unlock()
		}
	}()

	// Close the log stream when done
	defer log_stream.Close()

	client_gone := c.Writer.CloseNotify()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-client_gone:
			return false
		case logs, ok := <-log_chan:
			if !ok {
				// Channel is closed
				return false
			}

			w.Write(append(logs, byte('\n')))
		}
		return true
	})
}
