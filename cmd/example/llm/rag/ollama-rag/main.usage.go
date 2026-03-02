package main

import (
	"fmt"
)

func printUsage() {
	fmt.Println(`Usage:
  ollama-rag index --path <dir> [--namespace <ns>] [common flags]
  ollama-rag index --git-repo <repo> [--branch <name> | --tag <name>] [--log-limit <n>] [common flags]
  ollama-rag search --query <text> [--namespace <ns>] [retrieval flags] [common flags]
  ollama-rag ask --prompt <text> [--namespace <ns>] [retrieval flags] [common flags]

Common flags:
  --qdrant-url         Qdrant URL (default http://localhost:6333)
  --qdrant-collection  Qdrant collection (default default)
  --ollama-url         Ollama URL (default http://localhost:11434)
  --embed-model        Embedding model (default nomic-embed-text:latest)
  --model              Chat model for ask/query (default mistral:latest)

Retrieval flags:
  --docs               Number of docs to retrieve (default 5)
  --threshold          Score threshold [0..1] (default 0.4)

Notes:
  - Namespace filtering is metadata-based; omit --namespace for global search.
  - Git namespace is derived as <abs_repo_path>::<branch_or_tag_or_ref>.`)
}
