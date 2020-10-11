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

	fmt.Print(r.PromptMessage)

	text, _ := reader.ReadString('\n')
	return text
}
