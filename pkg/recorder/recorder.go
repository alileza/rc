package recorder

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/gordonklaus/portaudio"
)

// Config holds the configuration for the audio recorder
type Config struct {
	SampleRate    int
	ChunkDuration time.Duration
	BufferSize    int
	OutputDir     string
	NumChannels   int
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		SampleRate:    44100,
		ChunkDuration: time.Minute,
		BufferSize:    1024,
		OutputDir:     "recordings",
		NumChannels:   1, // Using mono by default
	}
}

// Recorder handles the audio recording process
type Recorder struct {
	config        *Config
	currentBuffer *audio.IntBuffer
	startTime     time.Time
	chunkNumber   int
	stream        *portaudio.Stream
}

// NewRecorder creates a new Recorder instance
func NewRecorder(config *Config) *Recorder {
	return &Recorder{
		config:      config,
		chunkNumber: 1,
	}
}

// Start begins the recording process
func (r *Recorder) Start() error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(r.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}

	// Initialize PortAudio
	if err := portaudio.Initialize(); err != nil {
		return fmt.Errorf("error initializing PortAudio: %v", err)
	}

	// Open default input stream
	stream, err := portaudio.OpenDefaultStream(
		r.config.NumChannels,
		0,
		float64(r.config.SampleRate),
		r.config.BufferSize,
		r.processAudio,
	)
	if err != nil {
		portaudio.Terminate()
		return fmt.Errorf("error opening stream: %v", err)
	}

	r.stream = stream

	// Start the stream
	if err := stream.Start(); err != nil {
		stream.Close()
		portaudio.Terminate()
		return fmt.Errorf("error starting stream: %v", err)
	}

	return nil
}

// Stop ends the recording process
func (r *Recorder) Stop() error {
	if r.stream != nil {
		r.stream.Stop()
		r.stream.Close()
	}
	portaudio.Terminate()
	return nil
}

// processAudio handles the incoming audio data
func (r *Recorder) processAudio(in []float32) {
	if r.currentBuffer == nil {
		// Initialize new buffer at the start of each chunk
		r.currentBuffer = &audio.IntBuffer{
			Format: &audio.Format{
				NumChannels: r.config.NumChannels,
				SampleRate:  r.config.SampleRate,
			},
			Data: make([]int, 0),
		}
		r.startTime = time.Now()
	}

	// Convert float32 samples to int
	for _, sample := range in {
		r.currentBuffer.Data = append(r.currentBuffer.Data, int(sample*32767))
	}

	// Check if we've reached the chunk duration
	if time.Since(r.startTime) >= r.config.ChunkDuration {
		// Save the current chunk
		filename := filepath.Join(r.config.OutputDir, fmt.Sprintf("chunk_%d.wav", r.chunkNumber))
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			return
		}

		// Create WAV encoder
		enc := wav.NewEncoder(file, r.config.SampleRate, 16, r.config.NumChannels, 1)

		// Write to WAV file
		if err := enc.Write(r.currentBuffer); err != nil {
			fmt.Printf("Error writing to WAV file: %v\n", err)
		}

		// Close the WAV file
		if err := enc.Close(); err != nil {
			fmt.Printf("Error closing WAV file: %v\n", err)
		}
		file.Close()

		fmt.Printf("Saved chunk %d to %s\n", r.chunkNumber, filename)
		r.chunkNumber++

		// Reset buffer for next chunk
		r.currentBuffer = nil
	}
}
