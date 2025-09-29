package main

import (
	"math"
)

func calculateRMS(buffer []float32) float32 {
	var sum float64
	for _, v := range buffer {
		sum += float64(v * v)
	}
	return float32(math.Sqrt(sum / float64(len(buffer))))
}
