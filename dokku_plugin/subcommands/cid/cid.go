package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/dokku/dokku/plugins/common"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: dokku haas:cid <app name>")
		os.Exit(1)
	}

	app_name := os.Args[2]

	if err := common.VerifyAppName(app_name); err != nil {
		fmt.Println("App not found")
		os.Exit(1)
	}

	file, err := os.Open(path.Join(os.Getenv("DOKKU_ROOT"), app_name, "CONTAINER.web.1"))
	if err != nil {
		fmt.Println("Error getting container ID: " + err.Error())
		os.Exit(1)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file: " + err.Error())
		os.Exit(1)
	}

	cid := strings.TrimSpace(string(content))

	fmt.Println(cid)
}
