package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/db"
)

func EnsureAuthedUser(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"status": "error", "message": "not_authed"})
		return
	}

	var user db.User
	result := db.DB.Joins("JOIN tokens ON tokens.user_id = users.id").Where("tokens.token = ?", token).First(&user)

	if result.Error != nil {
		c.AbortWithStatusJSON(401, gin.H{"status": "error", "message": "not_authed"})
		return
	}

	c.Set("user", user)

	c.Next()
}
