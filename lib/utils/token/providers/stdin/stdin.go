package stdin

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type StdInReader struct {
	PromptMessage string
}

func (r *StdInReader) Get() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Fprint(os.Stderr, r.PromptMessage)

	text, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("error reading from stdin: %v", err)
		return "", err
	}

	text = strings.TrimSuffix(text, "\n")
	if len(text) == 0 {
		return "", fmt.Errorf("empty input provided")
	}

	return text, nil
}
