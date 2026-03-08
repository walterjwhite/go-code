package runner

import (
	"errors"
	"io"
	"strings"
)

const (
	MaxIndentation = 1024
	MinIndentation = 0
)

type IndenterConf struct {
	Indentation int
	w           io.Writer
}

func Default(w io.Writer) (*IndenterConf, error) {
	return New(2, w)
}

func New(indentation int, w io.Writer) (*IndenterConf, error) {
	if indentation < MinIndentation {
		return nil, errors.New("indentation cannot be negative")
	}
	if indentation > MaxIndentation {
		return nil, errors.New("indentation exceeds maximum allowed value")
	}
	if w == nil {
		return nil, errors.New("writer cannot be nil")
	}
	return &IndenterConf{Indentation: indentation, w: w}, nil
}

func (i *IndenterConf) Write(p []byte) (n int, err error) {
	if i.w == nil {
		return 0, errors.New("writer is nil")
	}
	if i.Indentation < 0 || i.Indentation > MaxIndentation {
		return 0, errors.New("invalid indentation value")
	}

	prefix := strings.Repeat(" ", i.Indentation)

	prefixed := append([]byte(prefix), p...)
	return i.w.Write(prefixed)
}

func SafeNew(indentation int, w io.Writer) *IndenterConf {
	indenter, err := New(indentation, w)
	if err != nil {
		indenter, _ = Default(w)
	}
	return indenter
}
