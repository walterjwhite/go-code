package runner

import (
	"strings"
	"io"
)

type IndenterConf struct {
  Indentation int
  w io.Writer
}

func Default(w io.Writer) *IndenterConf {
	return New(2, w)
}

func New(indentation int, w io.Writer) *IndenterConf {
	return &IndenterConf{Indentation: indentation, w: w}
}

func (i *IndenterConf) Write(p []byte) (n int, err error) {
	return i.w.Write(append([]byte(strings.Repeat(" ", i.Indentation)), p...))
}
