package agent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRead(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "questions.txt")
	content := "What is Go?\nHow to test?\n"

	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	c := &Conf{
		QuestionFile: tmpFile,
	}

	c.read()

	expectedCount := 2
	if len(c.questions) != expectedCount {
		t.Errorf("expected %d questions, got %d", expectedCount, len(c.questions))
	}

	if len(c.questions) > 0 && c.questions[0] != "What is Go?" {
		t.Errorf("expected first question to be 'What is Go?', got %s", c.questions[0])
	}
}

