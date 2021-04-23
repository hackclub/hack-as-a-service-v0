package teams

import (
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/db"
)

func handlePOSTTeam(c *gin.Context) {
	var json struct {
		Name      string
		Automatic bool
		HNUserID  string
	}

	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid JSON"})
		return
	}

	// create in db
	team := db.Team{
		Name:      json.Name,
		Automatic: json.Automatic,
		HNUserID:  json.HNUserID,
	}
	result := db.DB.Create(&team)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "teamID": team.ID})
	}
}
