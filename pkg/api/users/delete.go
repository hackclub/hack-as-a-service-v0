package users

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
)

func handleDELETEUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid user ID"})
		return
	}

	result := db.DB.Delete(&db.User{}, id)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	}

	// Delete the user's team on their behalf too
	result = db.DB.Delete(&db.Team{}, "personal = TRUE AND id = ?", id)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
		return
	} else {
		c.JSON(200, gin.H{"status": "ok"})
	}
}
