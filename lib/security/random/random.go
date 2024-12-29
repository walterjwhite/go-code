package random

import (
	"math/rand"
)

func Of(deviation int) int {
	return rand.Intn(deviation)
}
