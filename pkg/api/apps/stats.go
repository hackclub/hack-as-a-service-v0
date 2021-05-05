package apps

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/api/util"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/hackclub/hack-as-a-service/pkg/irs"
)

func handleGETStats(c *gin.Context) {
	upgrader := util.MakeWebsocketUpgrader()

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

	ch := irs.CreateStatsOutput(app.ID)

	// Spin up a websocket connection
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	defer irs.RemoveStatsOutput(app.ID, ch)
	defer ws.Close()

	for {
		output := <-ch
		// log.Printf("Output = %+v\n", output)
		err := ws.WriteJSON(output)
		if err != nil {
			log.Printf("Error writing to ws: %+v\n", err)
			break
		}
	}
}
