package gh

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v35/github"
)

// Generate a github client to make requests against their REST API with
func genClient(installIDRaw string) (*github.Client, error) {
	// Getting app id
	appIDRaw, found := os.LookupEnv("GITHUB_APP_ID")
	if !found {
		return nil, errors.New("Failed to find github app ID env var")
	}

	// Parsing app ID to int
	appID, err := strconv.Atoi(appIDRaw)
	if err != nil {
		return nil, err
	}

	// Parsing install ID to int
	installID, err := strconv.Atoi(installIDRaw)
	if err != nil {
		return nil, err
	}

	// Creating client
	itr, err := ghinstallation.NewKeyFromFile(
		http.DefaultTransport,
		int64(appID),
		int64(installID),
		"hack-as-a-service.private-key.pem",
	)
	if err != nil {
		log.Printf("Failed to create github client key from ")
	}

	return github.NewClient(&http.Client{Transport: itr}), nil
}
