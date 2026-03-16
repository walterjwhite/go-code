package densecode

import (
	"fmt"
	"log"
)

func ExampleRenderer() {
	cfg := WithDefaults()
	cfg.ModuleSize = 10
	cfg.BitsPerModule = 3

	result, err := cfg.EncodeText("Hello, DenseCode!")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Rendered %d segments\n", len(result.Segments))
	for i, segment := range result.Segments {
		fmt.Printf("Segment %d: size=%d, modules=%d\n",
			i+1, segment.Code.size, segment.Code.ModuleSize)
	}
}

type CustomRenderer struct {
	Prefix string
}

func (r *CustomRenderer) Render(result *EncodeResult) error {
	fmt.Printf("%s: Rendering %d segments\n", r.Prefix, len(result.Segments))
	for i, segment := range result.Segments {
		fmt.Printf("%s: Segment %d - Size: %dx%d, Modules: %d\n",
			r.Prefix, i+1, segment.Code.size, segment.Code.size, segment.Code.ModuleSize)
	}
	return nil
}

func ExampleCustomRenderer() {
	cfg := WithDefaults()
	result, err := cfg.EncodeText("Custom renderer example")
	if err != nil {
		log.Fatal(err)
	}

	custom := &CustomRenderer{Prefix: "CUSTOM"}
	if err := custom.Render(result); err != nil {
		log.Fatal(err)
	}
}
