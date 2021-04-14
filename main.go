package main

import (
	"github.com/hackclub/hack-as-a-service/dokku"
)

func main() {
	dokku.RunCommand("help")
}
