package geo

import (
	"fmt"
	geol "github.com/kellydunn/golang-geo"
	"testing"
)

func TestDistance(t *testing.T) {
	p1 := geol.NewPoint(42.25, 120.2)
	p2 := geol.NewPoint(30.25, 112.2)

	d := Distance(p1, p2)

	fmt.Printf("great circle distance: %f\n", d)
}

func TestDistance2(t *testing.T) {
	p1 := geol.NewPoint(40.588742, -80.232017)
	p2 := geol.NewPoint(40.583878, -80.228884)

	d := Distance(p1, p2)

	fmt.Printf("great circle distance: %f\n", d)
}
