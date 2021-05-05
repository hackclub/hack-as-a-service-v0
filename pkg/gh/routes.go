package gh

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v35/github"
)

func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/webhook", handleWebhook)
}

func handleWebhook(c *gin.Context) {
	// Verify the payload and unload it
	payload, err := github.ValidatePayload(
		c.Request,
		[]byte(os.Getenv("GITHUB_WEBHOOK_SECRET")),
	)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid payload"})
		return
	}

	// Parse response
	event, err := github.ParseWebHook(github.WebHookType(c.Request), payload)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Invalid payload"})
		return
	}

	// Route type to appropriate handler
	switch e := event.(type) {
	case *github.PushEvent:
		handleDeployment(e)
	case *github.InstallationEvent:
		handleInstallEvent(e)
	}
}
