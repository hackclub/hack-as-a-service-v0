package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "haas",
		Short: "Hack Club's compute as a service",
		Long:  `cli to help you easily intergrate with`,
		Run: func(cmd *cobra.Command, args []string) {

			err := cmd.Help()

			if err != nil {
				fmt.Println("Failed to display help message")
			}

			os.Exit(0)

		},
	}

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
