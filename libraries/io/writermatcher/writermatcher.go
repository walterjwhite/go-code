package writermatcher

import (
	"io"
	"strings"
)

type WriterDelegate struct {
	Channel  chan *string
	Delegate io.Writer
	Matcher  WriterMatcher
}

type WriterMatcher interface {
	Check(p []byte)
}

func (writerDelegate WriterDelegate) Write(p []byte) (n int, err error) {
	writerDelegate.Matcher.check(p)

	if writerDelegate.Delegate != nil {
		return writerDelegate.Delegate.Write(p)
	}

	return len(p), nil
}

func onEachLine(s interface{}, p []byte, function func(s interface{}, line string)) {
	output := strings.Split(string(p), "\n")

	for _, line := range output {
		function(s, line)
	}
}
