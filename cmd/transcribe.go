package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

// TranscribeCommand is the command for transcribing audio files
var TranscribeCommand = &cli.Command{
	Name:  "transcribe",
	Usage: "Transcribe audio files using Whisper large-v3-turbo",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "input-dir",
			Value:   "recordings_maya_20250403",
			Usage:   "Directory containing WAV files to transcribe",
			Aliases: []string{"i"},
		},
		&cli.StringFlag{
			Name:    "output-dir",
			Value:   "transcripts",
			Usage:   "Directory to save transcriptions",
			Aliases: []string{"o"},
		},
	},
	Action: func(c *cli.Context) error {
		inputDir := c.String("input-dir")
		outputDir := c.String("output-dir")

		// Create output directory if it doesn't exist
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("error creating output directory: %v", err)
		}

		// Get all WAV files
		files, err := os.ReadDir(inputDir)
		if err != nil {
			return fmt.Errorf("error reading input directory: %v", err)
		}

		// Filter and sort WAV files
		var wavFiles []string
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".wav") {
				wavFiles = append(wavFiles, file.Name())
			}
		}

		// Process each file
		for _, wavFile := range wavFiles {
			fmt.Printf("\nProcessing %s...\n", wavFile)

			inputPath := filepath.Join(inputDir, wavFile)
			outputPath := filepath.Join(outputDir, strings.TrimSuffix(wavFile, ".wav")+"_transcript.txt")

			// Run Python script
			cmd := exec.Command("python3", "transcribe.py", inputPath, outputPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				fmt.Printf("Error transcribing %s: %v\n", wavFile, err)
				continue
			}

			fmt.Printf("Transcription saved to %s\n", outputPath)
		}

		return nil
	},
}
