package densecode

import "fmt"

type DecodeResult struct {
	Data           []byte
	IsMultiSegment bool
	TotalSegments  int
}

func (cfg *Configuration) DecodeFiles(result *EncodeResult) ([]byte, error) {
	if result == nil || len(result.Segments) == 0 {
		return nil, fmt.Errorf("no segments to decode")
	}

	if result.IsMultiSegment {
		segments := result.Segments
		return DecodeSegments(segments, cfg)
	}

	config := result.Segments[0].Code
	return config.Decode(config.ToMatrix())
}
