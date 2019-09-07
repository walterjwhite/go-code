package writermatcher

import (
	"io"
	"strings"
)

type ExceptionFilter struct {
	Channel    chan *string
	Patterns   []string
	IgnoreCase bool
}

func NewExceptionFilter(patterns []string, channel chan *string, writer io.Writer) WriterDelegate {
	return WriterDelegate{channel, writer, ExceptionFilter{Channel: channel, Patterns: patterns}}
}

func NewExceptionFilterIgnoreCase(patterns []string, channel chan *string, writer io.Writer) WriterDelegate {
	return WriterDelegate{channel, writer, ExceptionFilter{Channel: channel, Patterns: patterns, IgnoreCase: true}}
}

func (exceptionFilter ExceptionFilter) Check(p []byte) {
	onEachLine(exceptionFilter, p, checkLineForException)
}

func checkLineForException(s interface{}, line string) {
	field, ok := s.(ExceptionFilter)
	if !ok {
		panic("Error converting")
	}

	for _, pattern := range field.Patterns {
		if matches(line, pattern, field.IgnoreCase) {
			field.Channel <- &line
		}
	}
}

func matches(line string, pattern string, ignoreCase bool) bool {
	if ignoreCase {
		line = strings.ToUpper(line)
		pattern = strings.ToUpper(pattern)
	}

	return strings.Contains(line, pattern)
}
