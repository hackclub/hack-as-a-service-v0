package main

import (
	"encoding/json"
	"fmt"
	"os"

	dokku "github.com/dokku/dokku/plugins/common"
)

func main() {
	apps, err := dokku.DokkuApps()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	apps_json, err := json.Marshal(apps)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(string(apps_json))
}
