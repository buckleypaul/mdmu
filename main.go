package main

import (
	"os"

	"github.com/paulbuckley/mdmu/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
