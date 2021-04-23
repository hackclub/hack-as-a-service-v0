package oauth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hackclub/hack-as-a-service/db"
	"github.com/slack-go/slack"
	"gorm.io/gorm"
)

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

func SetupRoutes(r *gin.Engine) {
	r.GET("/login", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, generateAuthUrl())
	})
	r.GET("/oauth/login", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, generateAuthUrl())
	})

	r.GET("/oauth/code", func(c *gin.Context) {
		code := c.Query("code")
		resp, err := slack.GetOAuthV2Response(http.DefaultClient, os.Getenv("SLACK_CLIENT_ID"), os.Getenv("SLACK_CLIENT_SECRET"), code, getRedirectUri())
		if err != nil {
			c.String(500, err.Error())
			return
		}

		// token := resp.AuthedUser.AccessToken

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
			}

			db.DB.Create(&user)
		} else if result.Error != nil {
			// Something went wrong
			c.String(500, result.Error.Error())
			return
		}

		// Generate a token
		token := &db.Token{
			UserID: user.ID,
			Token:  generateToken(),
		}
		db.DB.Create(&token)

		c.SetCookie("token", token.Token, 2592000, "/", "", true, false)

		c.Redirect(http.StatusTemporaryRedirect, "/")
	})
}
