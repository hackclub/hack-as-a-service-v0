package auth

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/db"
)

func GetAuthedUser(c *gin.Context) (db.User, error) {
	token, err := c.Cookie("token")
	if err != nil {
		return db.User{}, err
	}

	log.Println(token)

	var user db.User
	result := db.DB.Raw("SELECT users.* FROM users INNER JOIN tokens ON tokens.user_id = users.id WHERE tokens.token = ? LIMIT 1", token).Scan(&user)

	if result.RowsAffected == 0 {
		return db.User{}, errors.New("token not found")
	}

	return user, nil
}
