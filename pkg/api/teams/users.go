package teams

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
)

func handlePOSTTeamUsers(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid app ID"})
		return
	}

	user := c.MustGet("user").(db.User)

	var json struct {
		User uint
	}
	err = c.BindJSON(&json)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid JSON"})
		return
	}

	var team db.Team
	result := db.DB.Joins("JOIN team_users ON teams.id = team_users.team_id").
		First(&team, "teams.id = ? AND team_users.user_id = ? AND NOT teams.personal", id, user.ID)
	if result.Error != nil {
		c.JSON(400, gin.H{"status": "error", "message": "team not found"})
		return
	}

	var invitedUser db.User
	result = db.DB.First(&invitedUser, "id = ?", json.User)
	if result.Error != nil {
		c.JSON(400, gin.H{"status": "error", "message": "user not found"})
		return
	}

	err = db.DB.Model(&team).Association("Users").Append(&invitedUser)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err})
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}
