package users

import (
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
)

func handleGETSearch(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(200, gin.H{"status": "ok", "message": "missing `q` parameter", "users": []interface{}{}})
		return
	}

	var users []db.User
	result := db.DB.Limit(10).Find(&users, "LOWER(name) LIKE CONCAT('%', LOWER(?), '%')", query)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "users": users})
}
