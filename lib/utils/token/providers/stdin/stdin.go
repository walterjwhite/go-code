package stdin

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	maxInputLength = 4096
)

type StdInReader struct {
	PromptMessage string
	Writer        io.Writer
}

func (r *StdInReader) Get() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	writer := r.Writer
	if writer == nil {
		writer = os.Stdout
	}

	_, err := fmt.Fprint(writer, r.PromptMessage)
	if err != nil {
		return "", err
	}

	text, err := reader.ReadString('\n')
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return "", fmt.Errorf("failed to read input")
		}
	}

	if len(text) > maxInputLength {
		return "", fmt.Errorf("input exceeds maximum length of %d characters", maxInputLength)
	}

	text = strings.ReplaceAll(text, "\x00", "") // Remove null bytes
	text = strings.TrimRight(text, "\r\n")      // Handle Windows (\r\n) and Unix (\n) line endings
	text = strings.TrimSpace(text)
	if len(text) == 0 {
		return "", fmt.Errorf("empty input provided")
	}

	return text, nil
}
