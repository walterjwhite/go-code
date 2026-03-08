package indexing

import "strings"

func ChunkText(input string, maxChars, overlap int) []string {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	if maxChars <= 0 {
		return nil
	}
	if overlap < 0 {
		overlap = 0
	}

	if len(input) <= maxChars {
		return []string{input}
	}
	if overlap >= maxChars {
		overlap = 0
	}

	chunks := make([]string, 0, len(input)/maxChars+1)
	start := 0
	for start < len(input) {
		end := min(start+maxChars, len(input))
		chunks = append(chunks, strings.TrimSpace(input[start:end]))
		if end == len(input) {
			break
		}
		start = end - overlap
	}

	return chunks
}
