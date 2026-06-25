package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/walterjwhite/go-code/lib/ai/rag"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

var (
	ragConfig = &rag.Config{}
	engine    *rag.Engine
)

func init() {
	application.Configure(ragConfig)
	logging.Error(ragConfig.New(application.Context), "init")
}

func main() {
	defer application.OnPanic()

	if len(os.Args) < 2 {
		fmt.Println("Usage: query <question>")
		return
	}

	question := strings.Join(os.Args[1:], " ")

	qctx, cancel := context.WithTimeout(application.Context, 5*time.Minute)
	defer cancel()

	answer(qctx, question)
}

func answer(ctx context.Context, question string) {
	fmt.Print("\n")
	docs, err := engine.Query(ctx, question, os.Stdout)
	logging.Error(err, "answer")

	fmt.Println()

	if len(docs) == 0 {
		return
	}

	fmt.Println("\nsources:")
	seen := map[string]bool{}
	for _, doc := range docs {
		src, _ := doc.Metadata["source"].(string)
		if src == "" || seen[src] {
			continue
		}
		seen[src] = true
		chunk, _ := doc.Metadata["chunk"].(int)
		fmt.Printf("   • %s  (chunk %d)\n", src, chunk)
	}
}
