package gh

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v35/github"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/hackclub/hack-as-a-service/pkg/util"
)

// Handle a POST request to add an installID to a certain team
func handleInstall(c *gin.Context) {
	// Get query parameters
	installID, err := util.ReqGetQuery("installation_id", c)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": err.Error()})
	}
	teamID, err := util.ReqGetQuery("team_id", c)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": err.Error()})
	}
	log.Println("Got install id of", installID)
	log.Println("Got team id of", installID)

	// Find the team with the teamID
	team := db.DB.Find(&db.Team{}, teamID)
	log.Printf("Team found: %#v\n", team)
}

// Handle an uninstall of the github app
// We basically just remove the installation ID associated with the team
func handleUninstall(c *github.InstallationEvent) {
	log.Println("Handled install action of", *c.Action)
}
