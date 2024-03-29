package teams

import (
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
)

func handlePOSTTeam(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	var json struct {
		Name      string
		Automatic bool
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
		Personal:  false,
		Users:     []*db.User{&user},
	}
	result := db.DB.Create(&team)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "team": team})
	}
}
