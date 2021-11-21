package client

import (
	"bufio"
	"encoding/json"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"os"
)

type FileFeed struct {
	Filename string
}

func (f *FileFeed) Fetch() []*Message {
	file, err := os.Open(f.Filename)
	logging.Panic(err)
	defer file.Close()

	messages := make([]*Message, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m := &Message{}
		logging.Panic(json.Unmarshal([]byte(scanner.Text()), m))
		logging.Panic(scanner.Err())

		messages = append(messages, m)
	}

	return messages
}
