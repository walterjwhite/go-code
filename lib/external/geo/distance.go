package geo

import (
	"github.com/kellydunn/golang-geo"
)

func Distance(p1, p2 *geo.Point) float64 {
	return p1.GreatCircleDistance(p2)
}
