package main

import (
	"sync"

	"github.com/gordonklaus/portaudio"
	"time"
)

type Transcriber struct {
	WhisperCLI   string
	ModelPath    string
	VADModelPath string

	TempDir string

	MessageHandlers []MessageHandler

	ChunkDuration time.Duration
	WallDuration  time.Duration

	state      *AudioState
	chunkQueue chan AudioChunk

	stream *portaudio.Stream
	wg     sync.WaitGroup

	outputFormat string
}

type AudioState struct {
	mu           sync.Mutex
	currentChunk []int16
	silenceCount int
	seq          int
}

type AudioChunk struct {
	Seq  int
	Data []int16
}
