package densecode

type DenseCodeRenderer interface {
	Render(result *EncodeResult) error
}
