# go-rag

A local RAG (Retrieval-Augmented Generation) application written in Go.  
Documents are embedded with **Ollama** and stored in a local **Qdrant** vector database.  
Questions are answered by an Ollama LLM grounded in your documents.

```
docs/
  paper.tex  notes.md  summary.txt
      │
      ▼  ingest
   Qdrant ──── embeddings (nomic-embed-text via Ollama)
      │
      ▼  query
   Ollama ──── answer (llama3.2) + source attribution
```

---

## Prerequisites

| Tool             | Install                             |
| ---------------- | ----------------------------------- |
| Docker + Compose | https://docs.docker.com/get-docker/ |
| Ollama           | https://ollama.com/download         |
| Go 1.22+         | https://go.dev/dl/                  |

---

## 1 — Start Qdrant

```bash
docker compose up -d
# Qdrant REST API is now at http://localhost:6333
# Optional web UI: http://localhost:6333/dashboard
```

---

## 2 — Pull Ollama models

```bash
# Embedding model (768-dimensional, fast, small)
ollama pull nomic-embed-text

# Chat model — pick any that fits your hardware:
ollama pull llama3.2        # 3 B params, good default
ollama pull llama3.1:8b     # 8 B params, higher quality
ollama pull mistral         # alternative
```

> **GPU note:** Ollama uses the GPU automatically if CUDA / Metal drivers are present.  
> The embedding model runs on CPU in seconds even without a GPU.

---

## 3 — Build

```bash
git clone <this-repo>
cd go-rag
go mod tidy      # downloads langchaingo, cobra, etc.
go build -o go-rag .
```

---

## 4 — Ingest documents

```bash
# Single file
./go-rag ingest ./docs/paper.tex

# Entire directory (recursive)
./go-rag ingest ./docs/

# Multiple paths at once
./go-rag ingest ./thesis/ ./notes/ ./README.md

# Custom chunk size (default: 800 chars, overlap: 100)
./go-rag ingest ./docs/ --chunk-size 600 --chunk-overlap 80
```

**Supported formats:** `.txt` `.text` `.md` `.mdx` `.markdown` `.tex`

You can re-ingest at any time; Qdrant deduplicates by point ID automatically.

---

## 5 — Query

```bash
./go-rag query
```

This starts an interactive session:

```
🔧 Connecting to Qdrant and Ollama…
  ✓ ready

Type a question and press Enter. Commands: /help  /quit
────────────────────────────────────────────────────────────
❓ What are the main contributions of the paper?

🤖 The paper introduces three main contributions: …

📚 Sources:
   • docs/paper.tex  (chunk 4)
   • docs/paper.tex  (chunk 7)
   • docs/notes.md   (chunk 1)

❓ /quit
Bye!
```

---

## 6 — Run the Pub/Sub worker

```bash
./go-rag serve \
  --gcp-project-id my-project \
  --gcp-credentials-file /path/to/service-account.json \
  --request-subscription rag-queries-sub \
  --response-topic rag-responses \
  --request-topic rag-queries
```

Request and response payloads are JSON:

```json
{ "ip_address": "203.0.113.10", "message": "What changed in the design doc?" }
```

The worker consumes a request, runs the RAG query, and publishes a response
with the same `ip_address` and the answer in `message`.

---

## Configuration flags

All flags are available on every subcommand:

| Flag               | Default                  | Description                |
| ------------------ | ------------------------ | -------------------------- |
| `--qdrant-url`     | `http://localhost:6333`  | Qdrant REST API URL        |
| `--ollama-url`     | `http://localhost:11434` | Ollama server URL          |
| `--collection`     | `documents`              | Qdrant collection name     |
| `--embed-model`    | `nomic-embed-text`       | Embedding model            |
| `--chat-model`     | `llama3.2`               | Chat / generation model    |
| `--ingest-timeout` | `5m`                     | Maximum runtime for ingest |

Ingest-only flags:

| Flag              | Default | Description                     |
| ----------------- | ------- | ------------------------------- |
| `--chunk-size`    | `800`   | Maximum characters per chunk    |
| `--chunk-overlap` | `100`   | Overlap between adjacent chunks |

Query-only flags:

| Flag         | Default | Description                   |
| ------------ | ------- | ----------------------------- |
| `--num-docs` | `5`     | Chunks retrieved per question |

Serve-only flags:

| Flag                           | Default | Description                                    |
| ------------------------------ | ------- | ---------------------------------------------- |
| `--request-topic`              | `""`    | Topic associated with the inbound subscription |
| `--request-subscription`       | `""`    | Subscription that receives query requests      |
| `--response-topic`             | `""`    | Topic to publish query responses to            |
| `--gcp-project-id`             | `""`    | Google Cloud project ID                        |
| `--gcp-credentials-file`       | `""`    | Service account credentials JSON path          |
| `--pubsub-encryption-key-file` | `""`    | Optional AES key file for lib/net/google       |
| `--pubsub-compress`            | `false` | Enable lib/net/google compression              |
| `--pubsub-serialize`           | `false` | Enable lib/net/google byte-slice serialization |
| `--query-timeout`              | `5m`    | Maximum runtime for each RAG query             |

### Example: different models or a remote Ollama

```bash
./go-rag ingest ./docs/ \
  --embed-model mxbai-embed-large \
  --ingest-timeout 15m \
  --collection research

./go-rag query \
  --chat-model llama3.1:8b \
  --collection research \
  --num-docs 8
```

---

## Multiple knowledge bases

Use `--collection` to maintain separate, isolated document sets:

```bash
./go-rag ingest ./thesis/   --collection thesis
./go-rag ingest ./receipts/ --collection finance

./go-rag query --collection thesis
./go-rag query --collection finance
```

---

## Resetting a collection

```bash
# Delete via Qdrant REST API
curl -X DELETE http://localhost:6333/collections/documents

# Or use the Qdrant dashboard
open http://localhost:6333/dashboard
```

---

## Architecture

```
cmd/
  root.go      — cobra root, shared flags
  ingest.go    — ingest subcommand
  query.go     — interactive query REPL
  serve.go     — long-running Pub/Sub-backed RAG worker

../../../lib/ai/
  document/
    loader.go  — file walker, format detection, Markdown/LaTeX cleaners
  rag/
    engine.go  — RAG engine: collection bootstrap, ingest pipeline,
                 similarity search, prompt assembly, streaming generation
    pubsub/
      service.go — Pub/Sub request consumer and response publisher
```

**Key design choices:**

- **Auto-detected vector size** — a probe embedding is generated on startup so the Qdrant collection is always created with the correct dimensionality, regardless of which embed model you choose.
- **Streaming answers** — tokens are written to stdout as they arrive, so you see the answer forming in real time.
- **Format-aware cleaning** — LaTeX and Markdown are stripped down to readable prose before chunking so the embedder works on semantic content, not markup noise.
- **No hard-coded model sizes** — swap any Ollama-compatible embedding or chat model via flags.
