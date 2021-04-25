package users

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
)

func handleGETAuthed(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	c.JSON(200, gin.H{"status": "ok", "user": user})
}

func handleGETUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid user ID"})
		return
	}

	var user db.User
	result := db.DB.First(&user, "id = ?", id)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error})
	} else {
		c.JSON(200, gin.H{"status": "ok", "user": user})
	}
}
