package main

import (
	"math"
	"os"

	"github.com/gordonklaus/portaudio"

	"github.com/walterjwhite/go-code/lib/application/logging"
)

func (t *Transcriber) init() error {
	err := portaudio.Initialize()
	if err != nil {
		return err
	}

	state := &AudioState{}
	chunkQueue := make(chan AudioChunk, 100)

	stream, err := openAudioStream(state, chunkQueue)
	if err != nil {
		_ = portaudio.Terminate()
		return err
	}

	t.state = state
	t.chunkQueue = chunkQueue

	t.stream = stream
	t.outputFormat = outputFormat

	return nil
}

func openAudioStream(state *AudioState, chunkQueue chan AudioChunk) (*portaudio.Stream, error) {
	dev, err := portaudio.DefaultInputDevice()
	if err != nil {
		return nil, err
	}
	params := portaudio.HighLatencyParameters(dev, nil)
	params.Input.Channels = channels
	params.Output.Channels = 0
	params.SampleRate = sampleRate
	params.FramesPerBuffer = framesPerBuffer
	maxSamples := int(maxChunkLength * sampleRate)
	silenceLimit := int(math.Ceil(float64(silenceDuration) * float64(sampleRate) / float64(framesPerBuffer)))
	return portaudio.OpenStream(params, func(in []float32) {
		state.mu.Lock()
		defer state.mu.Unlock()
		intBuf := make([]int16, len(in))
		for i, v := range in {
			v = float32(math.Max(-1, math.Min(1, float64(v))))
			intBuf[i] = int16(v * 32767)
		}
		state.currentChunk = append(state.currentChunk, intBuf...)
		rms := calculateRMS(in)
		if rms < silenceThresh {
			state.silenceCount++
		} else {
			state.silenceCount = 0
		}
		if state.silenceCount >= silenceLimit || len(state.currentChunk) >= maxSamples {
			if len(state.currentChunk) > 0 {
				chunkQueue <- AudioChunk{Seq: state.seq, Data: append([]int16(nil), state.currentChunk...)}
				state.seq++
				state.currentChunk = nil
				state.silenceCount = 0
			}
		}
	})
}

func (t *Transcriber) Start() error {
	err := t.stream.Start()
	if err != nil {
		return err
	}

	const numWorkers = 4
	for i := 0; i < numWorkers; i++ {
		t.wg.Add(1)
		go func() {
			defer t.wg.Done()
			t.processChunks()
		}()
	}

	return nil
}

func iclose(file *os.File) {
	_ = file.Close()
}

func remove(filename string) {
	_ = os.Remove(filename)
}

func (t *Transcriber) Process() {
	t.state.mu.Lock()
	if len(t.state.currentChunk) > 0 {
		t.chunkQueue <- AudioChunk{Seq: t.state.seq, Data: append([]int16(nil), t.state.currentChunk...)}
		t.state.seq++

		t.state.currentChunk = nil // Clear the current chunk after processing
		t.state.silenceCount = 0
	}
	t.state.mu.Unlock()
}

func (t *Transcriber) cleanup() {
	err := t.stream.Stop()
	logging.Warn(err, false, "Stop")

	close(t.chunkQueue)
	_ = portaudio.Terminate()

	t.wg.Wait()
}
