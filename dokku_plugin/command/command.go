package main

import (
	"fmt"
	"os"
	"strconv"
)

const HelpText = `Usage: dokku haas[:COMMAND]

Core plugin for HaaS

Additional commands:
    haas:apps				List apps in JSON format
    haas:cid <app name>			Get the container ID for an app
`

func main() {
	cmd := os.Args[1]

	switch cmd {
	case "help":
		fmt.Println("    haas, Core plugin for HaaS")
	case "haas":
		fallthrough
	case "haas:help":
		fmt.Println(HelpText)
	default:
		dokkuNotImplementExitCode, err := strconv.Atoi(os.Getenv("DOKKU_NOT_IMPLEMENTED_EXIT"))
		if err != nil {
			fmt.Println("failed to retrieve DOKKU_NOT_IMPLEMENTED_EXIT environment variable")
			dokkuNotImplementExitCode = 10
		}
		os.Exit(dokkuNotImplementExitCode)
	}
}
