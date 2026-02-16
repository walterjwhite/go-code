package pipe

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Flusher interface {
	Flush([]byte) error
}

type Reader struct {
	PipePath  string
	Threshold int
	Flusher   Flusher

	mu  sync.Mutex
	buf bytes.Buffer
}

func (r *Reader) Start() {
	go func() {
		for {
			f, err := os.OpenFile(r.PipePath, os.O_RDONLY, 0600)
			if err != nil {
				log.Error().Err(err).Str("path", r.PipePath).Msg("failed to open pipe for reading")
				time.Sleep(2 * time.Second)
				continue
			}

			reader := bufio.NewReader(f)

			for {
				chunk := make([]byte, 4096)
				n, err := reader.Read(chunk)
				if n > 0 {
					r.mu.Lock()
					r.buf.Write(chunk[:n])

					if r.Threshold > 0 && r.buf.Len() >= r.Threshold {
						if err := r.flushLocked(); err != nil {
							log.Error().Err(err).Msg("failed to flush buffer via flusher")
						}
					} else {
						r.mu.Unlock()
					}
				}

				if err != nil {
					if err == io.EOF {
						break
					}
					log.Error().Err(err).Msg("error reading from pipe")
					break
				}
			}

			if err := f.Close(); err != nil {
				log.Warn().Err(err).Msg("pipe.Reader.Start - closing pipe file")
			}

			time.Sleep(500 * time.Millisecond)
		}
	}()
}

func (r *Reader) Flush() error {
	r.mu.Lock()
	if r.buf.Len() == 0 {
		r.mu.Unlock()
		return nil
	}

	return r.flushLocked()
}

func (r *Reader) flushLocked() error {
	b := append([]byte(nil), r.buf.Bytes()...)
	r.buf.Reset()
	r.mu.Unlock()

	if r.Flusher == nil {
		return nil
	}

	if err := r.Flusher.Flush(b); err != nil {
		log.Error().Err(err).Int("bytes", len(b)).Msg("flush failed")
		return err
	}

	log.Info().Int("bytes", len(b)).Msg("flush succeeded")
	return nil
}

func (r *Reader) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.buf.Len()
}
