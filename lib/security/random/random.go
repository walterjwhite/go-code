package random

import (
	"math/rand"
	"time"
)

func Of(deviation int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(deviation)
}

