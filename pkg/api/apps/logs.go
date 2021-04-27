package apps

import (
	"io"
	"strings"

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
	if result.RowsAffected < 1 {
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

	// A channel to close the client stream when the Docker API disconnects
	log_stream_closed := make(chan bool)

	stdout_reader, stdout_writer := io.Pipe()
	stderr_reader, stderr_writer := io.Pipe()

	// Demultiplex stream in a goroutine
	go func() {
		defer stdout_reader.Close()
		defer stderr_reader.Close()

		defer func() { log_stream_closed <- true }()

		stdcopy.StdCopy(stdout_writer, stderr_writer, log_stream)
	}()

	// Make a channel to hold the log stream
	log_chan := make(chan []byte)

	// Listen for new logs in 2 seperate goroutines for stdout and stderr
	go func() {
		for {
			logs := make([]byte, 100)
			_, err := stdout_reader.Read(logs)
			if err != nil {
				return
			}

			log_chan <- logs
		}
	}()

	go func() {
		for {
			logs := make([]byte, 100)
			_, err := stderr_reader.Read(logs)
			if err != nil {
				return
			}

			log_chan <- logs
		}
	}()

	// Close the log stream when done
	defer log_stream.Close()

	client_gone := c.Writer.CloseNotify()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-client_gone:
			return false
		case <-log_stream_closed:
			return false
		case logs := <-log_chan:
			w.Write(logs)
		}
		return true
	})
}
