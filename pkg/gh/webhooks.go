package gh

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v35/github"
)

// Handle a webhook event of the following types:
// Push    -> Set deployment and checks
// Install -> Uninstall if that is the action type
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
		handlePush(e)
	case *github.InstallationEvent:
		if *e.Action == "uninstall" {
			handleUninstall(e)
		}
	}
}

// Handle a push event
func handlePush(e *github.PushEvent) {
	log.Println("Handled push to repo")
}
