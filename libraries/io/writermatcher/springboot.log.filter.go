package writermatcher

import (
	"fmt"
	"io"
	"strings"
)

// this only prints if the application log level is INFO or lower
const SpringBootApplicationStartedString = "Started Application in"

type SpringBootApplicationStartupNotifier struct {
	Channel  chan *string
	Notified bool
}

func New(channel chan *string, writer io.Writer) *WriterDelegate {
	return &WriterDelegate{Channel: channel, Delegate: writer, Matcher: SpringBootApplicationStartupNotifier{Channel: channel}}
}

func (n SpringBootApplicationStartupNotifier) Check(p []byte) {
	if n.Notified {
		return
	}

	onEachLine(n, p, check)
}

func check(s interface{}, line string) {
	field, ok := s.(SpringBootApplicationStartupNotifier)
	if !ok {
		panic(fmt.Sprintf("Error converting: %v\n", s))
	}

	if strings.Contains(line, SpringBootApplicationStartedString) {
		field.Channel <- &line
		field.Notified = true
	}
}
