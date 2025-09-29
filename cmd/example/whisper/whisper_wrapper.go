package main

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"os/exec"
	"path/filepath"
	"time"
)

func (t *Transcriber) processChunks() {
	for chunk := range t.chunkQueue {
		filename := filepath.Join(t.TempDir, fmt.Sprintf("chunk_%d.%s", time.Now().UnixNano(), outputFormat))

		log.Debug().Msgf("Processing chunk: %s (%d samples)", filename, len(chunk.Data))

		if err := generateWAV(filename, chunk.Data); err != nil {
			log.Warn().Msgf("Error saving chunk: %v", err)
			continue
		}

		cmd := exec.Command(
			t.WhisperCLI,
			"-m", t.ModelPath,
			"-np",
			"-nt",
			filename,
		)

		output, err := cmd.CombinedOutput()
		remove(filename)

		if err != nil {
			log.Warn().Msgf("Transcription failed: %v\n%s", err, output)
			continue
		}

		if err != nil {
			log.Warn().Msgf("Error reading transcription: %v", err)
		} else {
			if len(output) == 0 {
				log.Warn().Msg("nothing transcribed")
				continue
			}

			t.OnMessage(string(output))
		}
	}
}
