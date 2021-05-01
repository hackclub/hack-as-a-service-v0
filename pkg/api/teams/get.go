package teams

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
)

func handleGETTeam(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid billing account ID"})
		return
	}

	var team db.Team
	result := db.DB.Preload("Users").Preload("Apps").
		Joins("INNER JOIN team_users ON team_users.team_id = teams.id").
		First(&team, "teams.id = ? AND team_users.user_id = ?", id, user.ID)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error.Error()})
	} else {
		c.JSON(200, gin.H{"status": "ok", "team": team})
	}
}

func handleGETPersonalTeam(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	var team db.Team
	result := db.DB.Preload("Apps").Joins("JOIN team_users ON team_users.team_id = teams.id").First(&team, "team_users.user_id = ? AND teams.personal", user.ID)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error.Error()})
	} else {
		c.JSON(200, gin.H{"status": "ok", "team": team})
	}
}
