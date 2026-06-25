package document

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Format string

const (
	FormatText     Format = "text"
	FormatMarkdown Format = "markdown"
	FormatLatex    Format = "latex"
)

type Document struct {
	Content  string
	FilePath string
	Format   Format
}

func LoadPath(path string) ([]Document, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return loadDir(path)
	}
	doc, err := loadFile(path)
	if err != nil {
		return nil, err
	}
	return []Document{*doc}, nil
}

func loadDir(dir string) ([]Document, error) {
	var docs []Document
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if detectFormat(path) == "" {
			return nil
		}
		doc, err := loadFile(path)
		if err != nil {
			fmt.Printf("  warning: skipping %s: %v\n", path, err)
			return nil
		}
		docs = append(docs, *doc)
		return nil
	})
	return docs, err
}

func loadFile(path string) (*Document, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	format := detectFormat(path)
	return &Document{
		Content:  cleanContent(string(raw), format),
		FilePath: path,
		Format:   format,
	}, nil
}

func detectFormat(path string) Format {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".txt", ".text":
		return FormatText
	case ".md", ".mdx", ".markdown":
		return FormatMarkdown
	case ".tex":
		return FormatLatex
	}
	return ""
}

func cleanContent(content string, format Format) string {
	switch format {
	case FormatMarkdown:
		return cleanMarkdown(content)
	case FormatLatex:
		return cleanLatex(content)
	default:
		return content
	}
}

var (
	reMDFencedCode = regexp.MustCompile("(?s)```[^\\n]*\\n.*?```")
	reMDInlineCode = regexp.MustCompile("`[^`\\n]+`")
	reMDImage      = regexp.MustCompile(`!\[([^\]]*)\]\([^)]*\)`)
	reMDLink       = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	reMDHeader     = regexp.MustCompile(`(?m)^#{1,6}\s+`)
	reMDBold       = regexp.MustCompile(`\*{1,3}([^*\n]+)\*{1,3}`)
	reMDItalicUnd  = regexp.MustCompile(`_{1,2}([^_\n]+)_{1,2}`)
	reMDBlockquote = regexp.MustCompile(`(?m)^>\s+`)
	reMDHRule      = regexp.MustCompile(`(?m)^[-*_]{3,}\s*$`)
	reMDMultiBlank = regexp.MustCompile(`\n{3,}`)
)

func cleanMarkdown(s string) string {
	s = reMDFencedCode.ReplaceAllString(s, "[code block]")
	s = reMDInlineCode.ReplaceAllString(s, "")
	s = reMDImage.ReplaceAllString(s, "$1")
	s = reMDLink.ReplaceAllString(s, "$1")
	s = reMDHeader.ReplaceAllString(s, "")
	s = reMDBold.ReplaceAllString(s, "$1")
	s = reMDItalicUnd.ReplaceAllString(s, "$1")
	s = reMDBlockquote.ReplaceAllString(s, "")
	s = reMDHRule.ReplaceAllString(s, "")
	s = reMDMultiBlank.ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s)
}

var (
	reTexComment       = regexp.MustCompile(`(?m)%[^\n]*$`)
	reTexDisplayBrack  = regexp.MustCompile(`(?s)\\\[.*?\\\]`)
	reTexDisplayDollar = regexp.MustCompile(`(?s)\$\$.*?\$\$`)
	reTexInlineMath    = regexp.MustCompile(`\$[^$\n]+\$`)
	reTexFormatCmd     = regexp.MustCompile(`\\(?:text(?:bf|it|tt|rm|sf|sc|md)|emph|underline|uline)\{([^{}]*)\}`)
	reTexSectionCmd    = regexp.MustCompile(`\\(?:chapter|section|subsection|subsubsection|paragraph|subparagraph|caption|footnote|item)\{([^{}]*)\}`)
	reTexGenericArg    = regexp.MustCompile(`\\[a-zA-Z]+\{([^{}]*)\}`)
	reTexBeginEnd      = regexp.MustCompile(`\\(?:begin|end)\{[^}]*\}`)
	reTexCommand       = regexp.MustCompile(`\\[a-zA-Z]+\*?`)
	reTexBraces        = regexp.MustCompile(`[{}]`)
	reTexSpaces        = regexp.MustCompile(`[ \t]+`)
	reTexMultiBlank    = regexp.MustCompile(`\n{3,}`)
)

var nonContentEnvs = []string{
	"equation", "align", "align\\*", "eqnarray", "eqnarray\\*",
	"figure", "figure\\*", "table", "table\\*", "tabular", "tabular\\*",
	"lstlisting", "verbatim", "Verbatim", "tikzpicture", "pgfpicture",
	"algorithmic", "algorithm",
}

func cleanLatex(s string) string {
	s = reTexComment.ReplaceAllString(s, "")

	for _, env := range nonContentEnvs {
		re := regexp.MustCompile(`(?s)\\begin\{` + env + `\}.*?\\end\{` + env + `\}`)
		s = re.ReplaceAllString(s, "")
	}

	s = reTexDisplayBrack.ReplaceAllString(s, " [math] ")
	s = reTexDisplayDollar.ReplaceAllString(s, " [math] ")
	s = reTexInlineMath.ReplaceAllString(s, " [math] ")

	for range 5 {
		s = reTexFormatCmd.ReplaceAllString(s, "$1")
	}

	for range 5 {
		s = reTexSectionCmd.ReplaceAllString(s, "$1")
	}

	for range 4 {
		s = reTexGenericArg.ReplaceAllString(s, "$1")
	}

	s = reTexBeginEnd.ReplaceAllString(s, "")
	s = reTexCommand.ReplaceAllString(s, " ")
	s = reTexBraces.ReplaceAllString(s, " ")
	s = reTexSpaces.ReplaceAllString(s, " ")
	s = reTexMultiBlank.ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s)
}
