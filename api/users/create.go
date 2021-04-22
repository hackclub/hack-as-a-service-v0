package users

import (
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/db"
	"gorm.io/gorm"
)

func handlePOSTUser(c *gin.Context) {
	_db := c.MustGet("db").(*gorm.DB)
	var json struct {
		SlackUserID string
	}

	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid JSON"})
		return
	}

	// create in db
	user := db.User{SlackUserID: json.SlackUserID}
	result := _db.Create(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "userID": user.ID})
	}
}
