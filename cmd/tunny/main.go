package main

import (
	"os"

	"github.com/kayaramazan/tunny-client/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
