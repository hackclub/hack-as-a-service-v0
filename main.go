package main

import (
	"github.com/hackclub/hack-as-a-service/dokku"
	"fmt"
)

func main() {
	output, err := dokku.RunCommand("help")
	fmt.Printf(output)
	fmt.Printf("%s", err)
}
