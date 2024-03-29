package teams

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
)

func handleGETTeamApps(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid team ID"})
		return
	}

	var apps []db.App
	result := db.DB.
		Joins("INNER JOIN team_users ON team_users.team_id = apps.team_id").
		Where("apps.team_id = ? AND team_users.user_id = ?", uint(id), user.ID).Find(&apps)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "apps": apps})
	}
}
