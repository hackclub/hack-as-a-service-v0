package users

import (
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
)

func handleGETAuthed(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	c.JSON(200, gin.H{"status": "ok", "user": user})
}
