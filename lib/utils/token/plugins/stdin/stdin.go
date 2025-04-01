package stdin

import (
	"bufio"
	"fmt"
	"os"
)

type StdInReader struct {
	PromptMessage string
}

func (r *StdInReader) Get() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Fprint(os.Stderr, r.PromptMessage)

	text, _ := reader.ReadString('\n')
	return text
}
