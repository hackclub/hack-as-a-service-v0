package users

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
)

func handleGETSearch(c *gin.Context) {
	user := c.MustGet("user").(db.User)

	query := c.Query("q")
	if query == "" {
		c.JSON(200, gin.H{"status": "ok", "message": "missing `q` parameter", "users": []interface{}{}})
		return
	}

	limit := 10
	if query_limit := c.Query("limit"); query_limit != "" {
		if query_limit_int, err := strconv.Atoi(query_limit); err != nil {
			limit = query_limit_int
		}
	}

	exclude_self := c.Query("excludeSelf") != ""

	var users []db.User
	result := db.DB.Limit(limit).Find(&users, "LOWER(name) LIKE CONCAT('%', LOWER(?), '%') AND (? OR id != ?)", query, !exclude_self, user.ID)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "users": users})
}
