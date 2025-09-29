package main

import (
	"encoding/binary"
	"os"
)

func generateWAV(filename string, data []int16) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer iclose(file)

	header := make([]byte, 44)
	writeWAVHeader(header, len(data))
	if _, err := file.Write(header); err != nil {
		return err
	}

	buf := make([]byte, 2)
	for _, sample := range data {
		binary.LittleEndian.PutUint16(buf, uint16(sample))
		if _, err := file.Write(buf); err != nil {
			return err
		}
	}
	return nil
}

func writeWAVHeader(header []byte, dataSize int) {
	copy(header[0:], "RIFF")
	binary.LittleEndian.PutUint32(header[4:], uint32(36+dataSize*2))
	copy(header[8:], "WAVE")
	copy(header[12:], "fmt ")
	binary.LittleEndian.PutUint32(header[16:], 16)
	binary.LittleEndian.PutUint16(header[20:], 1) // PCM
	binary.LittleEndian.PutUint16(header[22:], channels)
	binary.LittleEndian.PutUint32(header[24:], sampleRate)
	binary.LittleEndian.PutUint32(header[28:], sampleRate*uint32(channels*2))
	binary.LittleEndian.PutUint16(header[32:], uint16(channels*2))
	binary.LittleEndian.PutUint16(header[34:], 16) // bits per sample
	copy(header[36:], "data")
	binary.LittleEndian.PutUint32(header[40:], uint32(dataSize*2))
}
