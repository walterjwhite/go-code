package indexing

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/schema"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/rag/git"
)

func validateGitRepoPath(repoPath string) error {
	if repoPath == "" {
		return fmt.Errorf("repository path cannot be empty")
	}
	if strings.Contains(repoPath, "..") {
		return fmt.Errorf("repository path contains invalid traversal sequence")
	}
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return fmt.Errorf("invalid repository path: %w", err)
	}
	realPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		return fmt.Errorf("invalid repository path: %w", err)
	}
	if _, err := os.Stat(realPath); os.IsNotExist(err) {
		return fmt.Errorf("repository path does not exist: %s", repoPath)
	}
	return nil
}

func sanitizeRef(ref string) error {
	if ref == "" {
		return fmt.Errorf("ref cannot be empty")
	}
	dangerousChars := []string{"..", ";", "|", "&", "$", "`", "(", ")", "<", ">", "\\", "\n", "\r"}
	for _, char := range dangerousChars {
		if strings.Contains(ref, char) {
			return fmt.Errorf("ref contains invalid characters")
		}
	}
	return nil
}

func IndexGitRepository(repoPath, branch, tag, namespace string, logLimit, maxFileBytes int) []schema.Document {
	if err := validateGitRepoPath(repoPath); err != nil {
		logging.Error(fmt.Errorf("invalid repository path: %w", err))
		return nil
	}

	if branch != "" {
		if err := sanitizeRef(branch); err != nil {
			logging.Error(fmt.Errorf("invalid branch name: %w", err))
			return nil
		}
	}
	if tag != "" {
		if err := sanitizeRef(tag); err != nil {
			logging.Error(fmt.Errorf("invalid tag name: %w", err))
			return nil
		}
	}

	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		logging.Error(fmt.Errorf("resolve absolute path: %w", err))
		return nil
	}

	realRepoPath, err := filepath.EvalSymlinks(absRepoPath)
	if err != nil {
		logging.Error(fmt.Errorf("resolve symlinks: %w", err))
		return nil
	}

	ref := git.ResolveRef(realRepoPath, branch, tag)
	if namespace == "" {
		namespace = fmt.Sprintf("%s::%s", realRepoPath, ref.Name)
	}

	docs := make([]schema.Document, 0, 1024)
	indexFilesDocs(&docs, realRepoPath, ref, namespace, maxFileBytes)
	indexLogDocs(&docs, realRepoPath, ref, namespace, logLimit)

	log.Info().Msgf("indexed %d git docs from repo=%s ref=%s namespace=%s", len(docs), realRepoPath, ref.Name, namespace)
	return docs
}

func indexFilesDocs(docs *[]schema.Document, repoPath string, ref git.Ref, namespace string, maxFileBytes int) {
	files, err := git.Lines(repoPath, "ls-tree", "-r", "--name-only", ref.Name)
	if err != nil {
		logging.Error(fmt.Errorf("list git files: %w", err))
	}

	for _, file := range files {
		content, readErr := git.ShowFile(repoPath, ref.Name, file)
		if readErr != nil {
			log.Warn().Err(readErr).Msgf("skip file %s", file)
			continue
		}
		if maxFileBytes > 0 && len(content) > maxFileBytes {
			continue
		}
		if git.IsBinary(content) {
			continue
		}

		chunks := ChunkText(string(content), 3000, 200)
		for chunkIdx, chunk := range chunks {
			*docs = append(*docs, schema.Document{
				PageContent: fmt.Sprintf("Repository: %s\nRef: %s\nPath: %s\n\n%s", repoPath, ref.Name, file, chunk),
				Metadata: map[string]any{
					"source":      "git",
					"doc_type":    "file",
					"repo_path":   repoPath,
					"file_path":   file,
					"ref_kind":    ref.Kind,
					"ref_name":    ref.Name,
					"branch_name": ref.BranchOrEmpty(),
					"tag_name":    ref.TagOrEmpty(),
					"namespace":   namespace,
					"chunk_index": chunkIdx,
				},
			})
		}
	}
}

func indexLogDocs(docs *[]schema.Document, repoPath string, ref git.Ref, namespace string, logLimit int) {
	if logLimit <= 0 {
		return
	}

	pretty := "%H%x1f%s%x1f%b%x1e"
	raw, err := git.Output(repoPath, "log", ref.Name, "--pretty=format:"+pretty, "-n", strconv.Itoa(logLimit))
	if err != nil {
		logging.Error(fmt.Errorf("read git log: %w", err))
	}

	entries := strings.SplitSeq(raw, "\x1e")
	for entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		parts := strings.Split(entry, "\x1f")
		if len(parts) < 3 {
			continue
		}

		hash := strings.TrimSpace(parts[0])
		subject := strings.TrimSpace(parts[1])
		body := strings.TrimSpace(parts[2])

		*docs = append(*docs, schema.Document{
			PageContent: fmt.Sprintf("Repository: %s\nRef: %s\nCommit: %s\nSubject: %s\n\n%s", repoPath, ref.Name, hash, subject, body),
			Metadata: map[string]any{
				"source":        "git",
				"doc_type":      "commit_log",
				"repo_path":     repoPath,
				"ref_kind":      ref.Kind,
				"ref_name":      ref.Name,
				"branch_name":   ref.BranchOrEmpty(),
				"tag_name":      ref.TagOrEmpty(),
				"log_message":   subject,
				"commit_hash":   hash,
				"namespace":     namespace,
				"content_scope": "log",
			},
		})
	}
}
