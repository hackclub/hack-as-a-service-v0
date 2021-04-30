package apps

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
)

func handleDELETEApp(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid app ID"})
		return
	}

	result := db.DB.Raw("DELETE FROM apps "+
		"INNER JOIN team_users ON team_users.team_id = apps.team_id "+
		"WHERE apps.id = ? AND team_users.user_id = ?", id, user.ID)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok"})
	}
}
