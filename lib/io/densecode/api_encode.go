package densecode

import (
	"fmt"
	"os"
)

type EncodeResult struct {
	data           []byte
	size           int
	Segments       []*Segment
	IsMultiSegment bool
}

func (cfg *Configuration) Encode(data []byte) (*EncodeResult, error) {
	cfg.initialize()

	if cfg.MaxSegmentSize > 0 && len(data) > cfg.MaxSegmentSize {
		segmentOpts := &SegmentOptions{
			Configuration:  cfg,
			MaxSegmentSize: cfg.MaxSegmentSize,
		}

		segments, err := EncodeSegments(data, segmentOpts)
		if err != nil {
			return nil, err
		}

		return &EncodeResult{
			Segments:       segments,
			IsMultiSegment: true,
		}, nil
	}

	maxCapacity := calculateMaxCapacity(cfg.ErrorLevel, cfg.BitsPerModule)
	if len(data) > maxCapacity {
		segmentOpts := &SegmentOptions{
			Configuration:  cfg,
			MaxSegmentSize: calculateOptimalSegmentSize(cfg.ErrorLevel, cfg.BitsPerModule),
		}

		segments, err := EncodeSegments(data, segmentOpts)
		if err != nil {
			return nil, err
		}

		return &EncodeResult{
			Segments:       segments,
			IsMultiSegment: true,
		}, nil
	}

	if err := cfg.encode(data); err != nil {
		return nil, err
	}

	segment := &Segment{
		Code:          cfg,
		SegmentIndex:  0,
		TotalSegments: 1,
	}

	return &EncodeResult{
		data:           cfg.data,
		size:           cfg.size,
		Segments:       []*Segment{segment},
		IsMultiSegment: false,
	}, nil
}

func (cfg *Configuration) EncodeText(text string) (*EncodeResult, error) {
	return cfg.Encode([]byte(text))
}

func (cfg *Configuration) EncodeFile(filePath string) (*EncodeResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return cfg.Encode(data)
}

func EncodeFiles(filePaths []string, cfg *Configuration) (*EncodeResult, error) {
	if len(filePaths) == 0 {
		return nil, fmt.Errorf("no files provided")
	}

	if len(filePaths) == 1 {
		return cfg.EncodeFile(filePaths[0])
	}

	tarData, err := createTarArchive(filePaths)
	if err != nil {
		return nil, err
	}

	return cfg.Encode(tarData)
}

func (cfg *Configuration) EncodeDirectory(dirPath string) (*EncodeResult, error) {
	tarData, err := createTarFromDirectory(dirPath)
	if err != nil {
		return nil, err
	}
	return cfg.Encode(tarData)
}
