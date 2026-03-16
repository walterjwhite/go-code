package png

import (
	"fmt"

	"github.com/walterjwhite/go-code/lib/io/densecode"
)

type PNGRenderer struct {
	BasePath string
}

func NewPNGRenderer(basePath string) *PNGRenderer {
	return &PNGRenderer{BasePath: basePath}
}

func (r *PNGRenderer) Render(result *densecode.EncodeResult) error {
	if result == nil || len(result.Segments) == 0 {
		return fmt.Errorf("no segments to render")
	}

	if result.IsMultiSegment {
		for i, segment := range result.Segments {
			filename := fmt.Sprintf("%s_%03d.png", r.BasePath, i+1)
			if err := r.renderSegment(segment.Code, filename); err != nil {
				return fmt.Errorf("failed to render segment %d: %w", i, err)
			}
		}
		return nil
	}

	return r.renderSegment(result.Segments[0].Code, r.BasePath)
}

func (r *PNGRenderer) renderSegment(code *densecode.Configuration, filename string) error {
	return code.RenderPNG(filename)
}
