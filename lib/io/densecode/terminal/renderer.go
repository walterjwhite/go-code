package terminal

import (
	"fmt"

	"github.com/walterjwhite/go-code/lib/io/densecode"
)

type TerminalRenderer struct {
	AdvanceKey bool
}

func NewTerminalRenderer(advanceKey bool) *TerminalRenderer {
	return &TerminalRenderer{AdvanceKey: advanceKey}
}

func (r *TerminalRenderer) Render(result *densecode.EncodeResult) error {
	if result == nil || len(result.Segments) == 0 {
		fmt.Println("No segments to display")
		return nil
	}

	if result.IsMultiSegment {
		fmt.Printf("Multi-segment DenseCode (%d segments):\n\n", len(result.Segments))
	}

	for i, segment := range result.Segments {
		if result.IsMultiSegment {
			fmt.Printf("Segment %d/%d:\n", i+1, len(result.Segments))
		}

		segment.Code.RenderTerminal()

		if r.AdvanceKey && i < len(result.Segments)-1 {
			fmt.Print("Press Enter to continue to next segment...")
			_, err := fmt.Scanln() // Wait for user input
			if err != nil {
				if err.Error() != "unexpected newline" {
					return fmt.Errorf("failed to read input: %w", err)
				}
			}
		}
	}

	return nil
}
