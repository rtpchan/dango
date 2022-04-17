package dango

import "math"

// general helper functions

// Tolerance check if a and b is close and within epsilon
func Tolerance(a, b, epsilon float64) bool {
	delta := math.Abs(a - b)
	return delta < epsilon

}

// DistanceTolerance check if two 2D positions are close enough
func DistanceTolerance(a, b Vector, epsilon float64) bool {
	delta := a.Distance(b)
	return delta < epsilon
}
