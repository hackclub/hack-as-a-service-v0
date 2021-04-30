package teams

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"gorm.io/gorm"
)

func handlePUTTeamUsers(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid app ID"})
		return
	}

	var json struct {
		Users []uint
	}
	err = c.BindJSON(&json)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid JSON"})
		return
	}

	var team db.Team
	result := db.DB.
		Joins("INNER JOIN team_users ON team_users.team_id = teams.id").
		First(&team, "teams.id = ? AND team_users.user_id = ?", id, user.ID)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
		return
	}
	team.Users = nil
	for _, user := range json.Users {
		team.Users = append(team.Users, &db.User{Model: gorm.Model{ID: user}})
	}

	result = db.DB.Save(&team)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok"})
	}
}
