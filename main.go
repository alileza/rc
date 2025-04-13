package main

import (
	"fmt"
	"os"

	"rc/cmd"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", err)
		}
	}

	app := &cli.App{
		Name:  "rc",
		Usage: "A tool for recording system audio output",
		Commands: []*cli.Command{
			cmd.RecordCommand,
			cmd.TranscribeCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
