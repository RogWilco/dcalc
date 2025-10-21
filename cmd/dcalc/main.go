package main

import (
	"os"

	"github.com/rogwilco/dcalc/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
