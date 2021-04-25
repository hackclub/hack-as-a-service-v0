package users

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
)

func handleGETUserTeams(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid user ID"})
		return
	}

	var teams []db.Team
	result := db.DB.
		Joins("INNER JOIN team_users ON team_users.team_id = teams.id").
		Where("team_users.user_id = ?", uint(id)).Find(&teams)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "teams": teams})
	}
}

func handleGETAuthedTeams(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	var teams []db.Team
	result := db.DB.
		Joins("INNER JOIN team_users ON team_users.team_id = teams.id").
		Where("team_users.user_id = ?", user.ID).Find(&teams)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "teams": teams})
	}
}
