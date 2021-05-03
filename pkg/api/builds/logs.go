package builds

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hackclub/hack-as-a-service/pkg/api/util"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/hackclub/hack-as-a-service/pkg/dokku"
)

func handleGETLogs(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid build ID"})
		return
	}

	var build db.Build
	result := db.DB.
		Joins("INNER JOIN apps ON apps.id = builds.app_id").
		Joins("INNER JOIN team_users ON team_users.team_id = apps.team_id").
		First(&build, "builds.id = ? AND team_users.user_id = ?", id, user.ID)
	if result.Error != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid build ID"})
		return
	}

	cmd, err := dokku.CreateOutput(build.ExecID)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": err.Error()})
	}

	// WebSockets
	upgrader := util.MakeWebsocketUpgrader()
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
	}

	defer func() {
		ws.Close()
		dokku.RemoveCommandOutput(build.ExecID, cmd)
	}()

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
