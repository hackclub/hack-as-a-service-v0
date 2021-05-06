package builds

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/hackclub/hack-as-a-service/pkg/dokku"
	"github.com/hackclub/hack-as-a-service/pkg/util"
)

func handleGETLogs(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	id, err := strconv.Atoi(c.Params.ByName("id"))
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
			err := ws.WriteJSON(gin.H{"Timestamp": time.Now().UnixNano(), "Stream": "stdout", "Output": line})
			switch err.(type) {
			case *websocket.CloseError:
				break loop
			}
		case line, ok := <-cmd.StderrChan:
			if !ok {
				continue
			}
			err := ws.WriteJSON(gin.H{"Timestamp": time.Now().UnixNano(), "Stream": "stderr", "Output": line})
			switch err.(type) {
			case *websocket.CloseError:
				break loop
			}
		case status := <-cmd.StatusChan:
			err = ws.WriteJSON(gin.H{"Timestamp": time.Now().UnixNano(), "Stream": "status", "Output": strconv.Itoa(status)})
			if err != nil {
				log.Println(err)
			}
			// no need to handle the error - we break anyways
			break loop
		}
	}
}
