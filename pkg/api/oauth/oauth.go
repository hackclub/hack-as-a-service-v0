package oauth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/slack-go/slack"
	"gorm.io/gorm"
)

func sweepOldTokens() error {
	result := db.DB.Unscoped().Delete(db.Token{}, "expires_at <= ?", time.Now())
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func getRedirectUri() string {
	redirect_uri := os.Getenv("SLACK_REDIRECT_URI")
	if redirect_uri == "" {
		redirect_uri = "https://haas.hackclub.com/oauth/code"
	}

	return redirect_uri
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func generateAuthUrl() string {
	oauth_url := url.URL{
		Scheme: "https",
		Host:   "slack.com",
		Path:   "/oauth/v2/authorize",
		RawQuery: url.Values{
			"redirect_uri": {getRedirectUri()},
			"user_scope":   {"identity.basic,identity.avatar"},
			"client_id":    {os.Getenv("SLACK_CLIENT_ID")},
		}.Encode(),
	}

	return oauth_url.String()
}

func SetupRoutes(r *gin.RouterGroup) {
	r.GET("/login", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, generateAuthUrl())
	})
	r.GET("/oauth/login", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, generateAuthUrl())
	})

	r.GET("/oauth/code", func(c *gin.Context) {
		code := c.Query("code")
		resp, err := slack.GetOAuthV2Response(
			http.DefaultClient,
			os.Getenv("SLACK_CLIENT_ID"),
			os.Getenv("SLACK_CLIENT_SECRET"),
			code,
			getRedirectUri(),
		)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		// Look for a user with that ID
		var user *db.User
		result := db.DB.Where("slack_user_id = ?", resp.AuthedUser.ID).First(&user)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create user
			client := slack.New(resp.AuthedUser.AccessToken)

			identity, err := client.GetUserIdentity()
			if err != nil {
				c.String(500, err.Error())
				return
			}

			user = &db.User{
				SlackUserID: resp.AuthedUser.ID,
				Name:        identity.User.Name,
				Avatar:      identity.User.Image512,
				Teams: []*db.Team{
					{
						Name:      "Personal team",
						Automatic: false,
						Personal:  true,
					},
				},
			}

			create_result := db.DB.Create(&user)
			if create_result.Error != nil {
				c.String(500, err.Error())
			}
		} else if result.Error != nil {
			// Something went wrong
			c.String(500, result.Error.Error())
			return
		}

		// Generate a token
		token := &db.Token{
			UserID:    user.ID,
			Token:     generateToken(),
			ExpiresAt: time.Now().Add(2592000 * time.Second),
		}
		create_result := db.DB.Create(&token)
		if create_result.Error != nil {
			c.String(500, result.Error.Error())
		}

		c.SetCookie("token", token.Token, 2592000, "/", "", true, true)

		c.Redirect(http.StatusTemporaryRedirect, "/")

		go sweepOldTokens()
	})

	r.GET("/logout", func(c *gin.Context) {
		defer c.Redirect(http.StatusTemporaryRedirect, "/")

		token, err := c.Cookie("token")
		if err != nil {
			return
		}

		// Clear token cookie
		c.SetCookie("token", "", 0, "/", "", true, false)

		// Revoke token
		result := db.DB.Delete(&db.Token{}, "token = ?", token)
		if result.Error != nil {
			log.Println(result.Error)
		}

		go sweepOldTokens()
	})
}
