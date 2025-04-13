package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"rc/pkg/recorder"

	"github.com/urfave/cli/v2"
)

// RecordCommand is the command for recording system audio
var RecordCommand = &cli.Command{
	Name:  "record",
	Usage: "Record system audio output in chunks",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output-dir",
			Value:   "recordings",
			Usage:   "Directory to save audio chunks",
			Aliases: []string{"o"},
		},
		&cli.DurationFlag{
			Name:    "chunk-duration",
			Value:   time.Minute,
			Usage:   "Duration of each audio chunk",
			Aliases: []string{"d"},
		},
		&cli.IntFlag{
			Name:    "sample-rate",
			Value:   44100,
			Usage:   "Audio sample rate",
			Aliases: []string{"r"},
		},
		&cli.IntFlag{
			Name:    "channels",
			Value:   1,
			Usage:   "Number of audio channels (1 for mono, 2 for stereo)",
			Aliases: []string{"c"},
		},
	},
	Action: func(c *cli.Context) error {
		// Create recorder configuration
		config := &recorder.Config{
			SampleRate:    c.Int("sample-rate"),
			ChunkDuration: c.Duration("chunk-duration"),
			BufferSize:    1024,
			OutputDir:     c.String("output-dir"),
			NumChannels:   c.Int("channels"),
		}

		// Create new recorder
		r := recorder.NewRecorder(config)

		// Handle interrupt signal
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)

		// Start recording
		fmt.Println("Starting audio recording... Press Ctrl+C to stop")
		fmt.Println("Make sure BlackHole is selected as your system output device in System Preferences > Sound > Output")

		if err := r.Start(); err != nil {
			return fmt.Errorf("error starting recorder: %v", err)
		}

		// Wait for interrupt signal
		<-sig
		fmt.Println("\nStopping recording...")
		return r.Stop()
	},
}
