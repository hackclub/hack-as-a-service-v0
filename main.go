package main

import (
	"fmt"

	"github.com/hackclub/hack-as-a-service/dokku"
)

func main() {
	output, err := dokku.RunCommand("help")
	fmt.Printf("%s", output)
	fmt.Printf("%s", err)
}
