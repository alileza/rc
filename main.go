package main

import (
	"fmt"
	"os"

	"rc/cmd"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "rc",
		Usage: "A tool for recording system audio output",
		Commands: []*cli.Command{
			cmd.RecordCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
