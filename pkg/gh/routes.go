package gh

import (
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v35/github"
)

func SetupRoutes(r *gin.RouterGroup) {
	r.POST("/webhook", handleWebhook)
}

func handleWebhook(c *gin.Context) {
	// Verifying the webhook
	// shaHeader := c.GetHeader("X-Hub-Signature-256")
	// log.Println("Header:", shaHeader)
	// if os.Getenv("GITHUB_WEBHOOK_SECRET_SHA") != shaHeader {
	// 	c.JSON(403, gin.H{"status": "error", "message": "Invalid sha256 header"})
	// 	return
	// }

	// Get raw data
	payload, err := c.GetRawData()
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
