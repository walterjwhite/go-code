package codesearch

import (
	"bytes"
	"fmt"
	"github.com/walterjwhite/go-application/libraries/logging"
	"io"
)

type Match struct {
	Filename   string
	LineNumber int
	Matched    []byte
}

var nl = []byte{'\n'}

func (s *SearchInstance) reader(r io.Reader, name string) {
	if s.buf == nil {
		s.buf = make([]byte, 1<<20)
	}
	var (
		buf = s.buf[:0]

		lineno = 1

		beginText = true
		endText   = false
	)

	for {
		n, err := io.ReadFull(r, buf[len(buf):cap(buf)])
		buf = buf[:len(buf)+n]
		end := len(buf)
		if err == nil {
			i := bytes.LastIndex(buf, nl)
			if i >= 0 {
				end = i + 1
			}
		} else {
			endText = true
		}
		chunkStart := 0
		for chunkStart < end {
			m1 := s.regexp.Match(buf[chunkStart:end], beginText, endText) + chunkStart
			beginText = false
			if m1 < chunkStart {
				break
			}

			lineStart := bytes.LastIndex(buf[chunkStart:m1], nl) + 1 + chunkStart
			lineEnd := m1 + 1
			if lineEnd > end {
				lineEnd = end
			}

			lineno += countNL(buf[chunkStart:lineStart])

			line := buf[lineStart:lineEnd]
			// nl = []byte{''}
			nl = nil
			if len(line) == 0 || line[len(line)-1] != '\n' {
				nl = []byte{'\n'}
			}

			// TODO: return new instance of match
			s.SearchOutputProcessor.OnMatch(&Match{Filename: name, LineNumber: lineno, Matched: line})

			lineno++

			chunkStart = lineEnd
		}
		if err == nil {
			lineno += countNL(buf[chunkStart:end])
		}
		n = copy(buf, buf[end:])
		buf = buf[:n]
		if len(buf) == 0 && err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				logging.Warn(err, false, fmt.Sprintf("Other unexpected error while searching file: %v", name))
			}
			break
		}
	}
}

func countNL(b []byte) int {
	n := 0
	for {
		i := bytes.IndexByte(b, '\n')
		if i < 0 {
			break
		}
		n++
		b = b[i+1:]
	}
	return n
}
