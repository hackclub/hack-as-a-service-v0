package gh

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v35/github"
)

// Generate a github client to make requests against their REST API with
func genClient(installIDRaw string) *github.Client {
	// Getting app id
	appIDRaw, found := os.LookupEnv("GITHUB_APP_ID")
	if !found {
		log.Fatalln("Failed to load app id for github bot. Is it stored?")
	}

	// Parsing app ID to int
	appID, err := strconv.Atoi(appIDRaw)
	if err != nil {
		log.Println(err)
		log.Fatalf("Failed to convert %q to int for github bot app id.\n", appIDRaw)
	}

	// Parsing install ID to int
	installID, err := strconv.Atoi(installIDRaw)
	if err != nil {
		log.Printf("Failed to convert %q to int for github bot installation id.\n", installIDRaw)
		log.Println(err)
	}

	// Creating client
	tr := http.DefaultTransport
	itr, err := ghinstallation.NewKeyFromFile(
		tr,
		int64(appID),
		int64(installID),
		"hack-as-a-service.private-key.pem",
	)
	if err != nil {
		log.Printf("Failed to create github client key from ")
	}

	return github.NewClient(&http.Client{Transport: itr})
}
