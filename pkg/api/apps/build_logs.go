package apps

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hackclub/hack-as-a-service/pkg/api/util"
	"github.com/hackclub/hack-as-a-service/pkg/dokku"
)

func handleGETBuildLogs(c *gin.Context) {
	_, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid app ID"})
		return
	}

	buildId, err := strconv.Atoi(c.Query("build_id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid build ID"})
		return
	}

	cmd, ok := dokku.GetRunningCommand(buildId)
	if !ok {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid build ID"})
	}

	// WebSockets
	upgrader := util.MakeWebsocketUpgrader()
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
	}

	defer ws.Close()

loop:
	for {
		select {
		case line, ok := <-cmd.StdoutChan:
			if !ok {
				continue
			}
			err := ws.WriteJSON(gin.H{"stdout": line})
			switch err.(type) {
			case *websocket.CloseError:
				break loop
			}
		case line, ok := <-cmd.StderrChan:
			if !ok {
				continue
			}
			err := ws.WriteJSON(gin.H{"stderr": line})
			switch err.(type) {
			case *websocket.CloseError:
				break loop
			}
		case status := <-cmd.StatusChan:
			ws.WriteJSON(gin.H{"status": status})
			// no need to handle the error - we break anyways
			break loop
		}
	}
}
