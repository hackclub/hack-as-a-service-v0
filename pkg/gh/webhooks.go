package gh

import (
	"log"

	"github.com/google/go-github/v35/github"
)

func handleDeployment(e *github.PushEvent) {
	log.Println("Handled push to repo")
}

func handleInstallEvent(c *github.InstallationEvent) {
	log.Println("Handled install action of " + *c.Action)
}
