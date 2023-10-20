package writermatcher

import (
	"fmt"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"io"
	"strings"
)

type LineMatcher struct {
	Channel chan *string

	Line       string
	Notified   bool
	NotifyOnce bool
}

func (n LineMatcher) Check(p []byte) {
	if n.Notified {
		return
	}

	onEachLine(n, p, check)
}

type ConversionError struct {
	Interface interface{}
}

func (ce *ConversionError) Error() string {
	return fmt.Sprintf("Error converting: %v\n", ce.Interface)
}

func check(s interface{}, line string) {
	field, ok := s.(LineMatcher)
	if !ok {
		logging.Panic(&ConversionError{Interface: s})
	}

	if strings.Contains(line, field.Line) {
		if !field.Notified || !field.NotifyOnce {
			field.Channel <- &line
		}

		field.Notified = true
	}
}

func NewLineMatcher(channel chan *string, writer io.Writer, line string, notifyOnce bool) *WriterDelegate {
	return &WriterDelegate{Channel: channel, Delegate: writer, Matcher: LineMatcher{Channel: channel, NotifyOnce: notifyOnce, Line: line}}
}

func NewSpringBootApplicationStartupMatcher(channel chan *string, writer io.Writer) *WriterDelegate {
	return NewLineMatcher(channel, writer, "Started Application in", true)
}

func NewNPMStartupMatcher(channel chan *string, writer io.Writer) *WriterDelegate {
	return NewLineMatcher(channel, writer, "webpack: Compiled successfully.", true)
}
