package audit

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"strings"

	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"

	"encoding/csv"
	"io"
)

func Run(ctx context.Context, scriptFile string, label string) {
	file, err := os.Open(scriptFile)
	logging.Panic(err)

	defer file.Close()

	r := csv.NewReader(file)

	i := 0
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}

		logging.Panic(err)

		if isComment(row[0]) {
			continue
		}

		childLabel := filepath.Join(label, fmt.Sprintf("%v.%v", i, row[0]))

		cmd := runner.Prepare(ctx, row[0], row[1:]...)

		Execute(cmd, childLabel)

		i++
	}
}

func isComment(command string) bool {
	return strings.HasPrefix(command, "#")
}
