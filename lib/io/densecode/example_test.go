package densecode

import (
	"fmt"
	"os"
)

func ExampleConfiguration_EncodeText() {
	cfg := &Configuration{}
	result, err := cfg.EncodeText("Hello, DenseCode!")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	err = result.Segments[0].Code.RenderPNG("output.png")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	matrix := result.Segments[0].Code.ToMatrix()
	data, err := cfg.Decode(matrix)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Decoded: %s\n", string(data))
}

func ExampleConfiguration_EncodeFile() {
	err := os.WriteFile("test.txt", []byte("File content"), 0644)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer func() { _ = os.Remove("test.txt") }()

	cfg := &Configuration{}
	result, err := cfg.EncodeFile("test.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	err = result.Segments[0].Code.RenderPNG("file_output.png")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("File encoded successfully")
}

func ExampleEncodeFiles() {
	files := []string{"file1.txt", "file2.txt"}
	for i, name := range files {
		err := os.WriteFile(name, fmt.Appendf(nil, "Content of file %d", i+1), 0644)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		defer func(n string) { _ = os.Remove(n) }(name)
	}

	result, err := EncodeFiles(files, &Configuration{})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if result.IsMultiSegment {
		for i, segment := range result.Segments {
			filename := fmt.Sprintf("multi_file_output_%03d.png", i+1)
			err = segment.Code.RenderPNG(filename)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		}
	} else {
		err = result.Segments[0].Code.RenderPNG("multi_file_output.png")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	}

	fmt.Printf("Encoded %d files into %d segment(s)\n", len(files), len(result.Segments))
}

func ExampleConfiguration_EncodeDirectory() {
	err := os.MkdirAll("testdir", 0755)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer func() { _ = os.RemoveAll("testdir") }()

	err = os.WriteFile("testdir/file1.txt", []byte("File 1"), 0644)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	err = os.WriteFile("testdir/file2.txt", []byte("File 2"), 0644)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	cfg := &Configuration{}
	result, err := cfg.EncodeDirectory("testdir")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	err = result.Segments[0].Code.RenderPNG("directory_output.png")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Directory encoded into %d segment(s)\n", len(result.Segments))
}

func ExampleEncode_withOptions() {
	data := []byte("Data with custom options")

	cfg := &Configuration{
		ErrorLevel: 2,          // Higher error correction
		ModuleSize: 8,          // Smaller modules
		Profile:    "balanced", // 3 bits per module
		Compressor: nil,        // Optional: add compressor
		Encryptor:  nil,        // Optional: add encryptor
	}

	result, err := cfg.Encode(data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	err = result.Segments[0].Code.RenderPNG("custom_output.png")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Encoded with %d segment(s)\n", len(result.Segments))
}

func ExampleEncode_largeData() {
	data := make([]byte, 50*1024) // 50KB
	for i := range data {
		data[i] = byte(i % 256)
	}

	cfg := &Configuration{
		MaxSegmentSize: 20 * 1024, // 20KB per segment
	}

	result, err := cfg.Encode(data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if result.IsMultiSegment {
		for i, segment := range result.Segments {
			filename := fmt.Sprintf("large_data_segment_%03d.png", i+1)
			err = segment.Code.RenderPNG(filename)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		}
	} else {
		err = result.Segments[0].Code.RenderPNG("large_data_segment.png")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	}

	fmt.Printf("Large data encoded into %d segments\n", len(result.Segments))
}

func ExampleConfiguration_RenderTerminal() {
	cfg := &Configuration{}
	result, _ := cfg.EncodeText("Hi")

	_ = result.Segments[0] // Use the result

	fmt.Println("Done")
}

func ExampleDecode() {
	original := "Secret message"

	cfg := &Configuration{}
	result, _ := cfg.EncodeText(original)
	matrix := result.Segments[0].Code.ToMatrix()

	data, err := cfg.Decode(matrix)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Decoded: %s\n", string(data))
}
