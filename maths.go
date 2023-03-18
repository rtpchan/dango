package dango

import (
	"errors"
	"math"

	"golang.org/x/image/math/f64"
)

// general helper functions

// Tolerance check if a and b is close and within epsilon
func Tolerance(a, b, epsilon float64) bool {
	delta := math.Abs(a - b)
	return delta < epsilon

}

// DistanceTolerance check if two 2D positions are close enough
func DistanceTolerance(a, b f64.Vec2, epsilon float64) bool {
	// delta := a.Distance(b)
	delta := math.Sqrt(math.Pow(a[0]-b[0], 2) + math.Pow(a[1]-b[1], 2))
	return delta < epsilon
}

// SegmentsIntersect find intersection point between line pt1 to pt2 and pt3 and pt4
// return error if no intersection
func SegmentsIntersect(x1, y1, x2, y2, x3, y3, x4, y4 float64) (float64, float64, error) {
	t := ((x1-x3)*(y3-y4) - (y1-y3)*(x3-x4)) /
		((x1-x2)*(y3-y4) - (y1-y2)*(x3-x4))
	u := ((x1-x3)*(y1-y2) - (y1-y3)*(x1-x2)) /
		((x1-x2)*(y3-y4) - (y1-y2)*(x3-x4))
	if t < 0 || t > 1 || u < 0 || u > 1 {
		return 0, 0, errors.New("lines do not intersect")
	}
	return x1 + t*(x2-x1), y1 + t*(y2-y1), nil
}
