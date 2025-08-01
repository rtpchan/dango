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

func EqualFloat(a, b, tolerance float64) bool {
	// tolerance := 1e-3
	delta := a - b
	return math.Abs(delta) < tolerance
}

func MatrixMultiplication(a, b []float64) []float64 {
	return []float64{
		a[0]*b[0] + a[1]*b[4] + a[2]*b[8] + a[3]*b[12],
		a[0]*b[1] + a[1]*b[5] + a[2]*b[9] + a[3]*b[13],
		a[0]*b[2] + a[1]*b[6] + a[2]*b[10] + a[3]*b[14],
		a[0]*b[3] + a[1]*b[7] + a[2]*b[11] + a[3]*b[15],

		a[4]*b[0] + a[5]*b[4] + a[6]*b[8] + a[7]*b[12],
		a[4]*b[1] + a[5]*b[5] + a[6]*b[9] + a[7]*b[13],
		a[4]*b[2] + a[5]*b[6] + a[6]*b[10] + a[7]*b[14],
		a[4]*b[3] + a[5]*b[7] + a[6]*b[11] + a[7]*b[15],

		a[8]*b[0] + a[9]*b[4] + a[10]*b[8] + a[11]*b[12],
		a[8]*b[1] + a[9]*b[5] + a[10]*b[9] + a[11]*b[13],
		a[8]*b[2] + a[9]*b[6] + a[10]*b[10] + a[11]*b[14],
		a[8]*b[3] + a[9]*b[7] + a[10]*b[11] + a[11]*b[15],

		a[12]*b[0] + a[13]*b[4] + a[14]*b[8] + a[15]*b[12],
		a[12]*b[1] + a[13]*b[5] + a[14]*b[9] + a[15]*b[13],
		a[12]*b[2] + a[13]*b[6] + a[14]*b[10] + a[15]*b[14],
		a[12]*b[3] + a[13]*b[7] + a[14]*b[11] + a[15]*b[15],
	}
}

func MatrixVectorMultiplication(a, b []float64) []float64 {
	return []float64{
		a[0]*b[0] + a[1]*b[1] + a[2]*b[2] + a[3]*b[3],
		a[4]*b[0] + a[5]*b[1] + a[6]*b[2] + a[7]*b[3],
		a[8]*b[0] + a[9]*b[1] + a[10]*b[2] + a[11]*b[3],
		a[12]*b[0] + a[13]*b[1] + a[14]*b[2] + a[15]*b[3],
	}
}

// InvertMatrix computes the inverse of a 4x4 matrix.
func InvertMatrix(m []float64) ([]float64, bool) {
	inv := make([]float64, 16)
	inv[0] = m[5]*m[10]*m[15] - m[5]*m[11]*m[14] - m[9]*m[6]*m[15] + m[9]*m[7]*m[14] + m[13]*m[6]*m[11] - m[13]*m[7]*m[10]
	inv[4] = -m[4]*m[10]*m[15] + m[4]*m[11]*m[14] + m[8]*m[6]*m[15] - m[8]*m[7]*m[14] - m[12]*m[6]*m[11] + m[12]*m[7]*m[10]
	inv[8] = m[4]*m[9]*m[15] - m[4]*m[11]*m[13] - m[8]*m[5]*m[15] + m[8]*m[7]*m[13] + m[12]*m[5]*m[11] - m[12]*m[7]*m[9]
	inv[12] = -m[4]*m[9]*m[14] + m[4]*m[10]*m[13] + m[8]*m[5]*m[14] - m[8]*m[6]*m[13] - m[12]*m[5]*m[10] + m[12]*m[6]*m[9]
	inv[1] = -m[1]*m[10]*m[15] + m[1]*m[11]*m[14] + m[9]*m[2]*m[15] - m[9]*m[3]*m[14] - m[13]*m[2]*m[11] + m[13]*m[3]*m[10]
	inv[5] = m[0]*m[10]*m[15] - m[0]*m[11]*m[14] - m[8]*m[2]*m[15] + m[8]*m[3]*m[14] + m[12]*m[2]*m[11] - m[12]*m[3]*m[10]
	inv[9] = -m[0]*m[9]*m[15] + m[0]*m[11]*m[13] + m[8]*m[1]*m[15] - m[8]*m[3]*m[13] - m[12]*m[1]*m[11] + m[12]*m[3]*m[9]
	inv[13] = m[0]*m[9]*m[14] - m[0]*m[10]*m[13] - m[8]*m[1]*m[14] + m[8]*m[2]*m[13] + m[12]*m[1]*m[10] - m[12]*m[2]*m[9]
	inv[2] = m[1]*m[6]*m[15] - m[1]*m[7]*m[14] - m[5]*m[2]*m[15] + m[5]*m[3]*m[14] + m[13]*m[2]*m[7] - m[13]*m[3]*m[6]
	inv[6] = -m[0]*m[6]*m[15] + m[0]*m[7]*m[14] + m[4]*m[2]*m[15] - m[4]*m[3]*m[14] - m[12]*m[2]*m[7] + m[12]*m[3]*m[6]
	inv[10] = m[0]*m[5]*m[15] - m[0]*m[7]*m[13] - m[4]*m[1]*m[15] + m[4]*m[3]*m[13] + m[12]*m[1]*m[7] - m[12]*m[3]*m[5]
	inv[14] = -m[0]*m[5]*m[14] + m[0]*m[6]*m[13] + m[4]*m[1]*m[14] - m[4]*m[2]*m[13] - m[12]*m[1]*m[6] + m[12]*m[2]*m[5]
	inv[3] = -m[1]*m[6]*m[11] + m[1]*m[7]*m[10] + m[5]*m[2]*m[11] - m[5]*m[3]*m[10] - m[9]*m[2]*m[7] + m[9]*m[3]*m[6]
	inv[7] = m[0]*m[6]*m[11] - m[0]*m[7]*m[10] - m[4]*m[2]*m[11] + m[4]*m[3]*m[10] + m[8]*m[2]*m[7] - m[8]*m[3]*m[6]
	inv[11] = -m[0]*m[5]*m[11] + m[0]*m[7]*m[9] + m[4]*m[1]*m[11] - m[4]*m[3]*m[9] - m[8]*m[1]*m[7] + m[8]*m[3]*m[5]
	inv[15] = m[0]*m[5]*m[10] - m[0]*m[6]*m[9] - m[4]*m[1]*m[10] + m[4]*m[2]*m[9] + m[8]*m[1]*m[6] - m[8]*m[2]*m[5]

	det := m[0]*inv[0] + m[1]*inv[4] + m[2]*inv[8] + m[3]*inv[12]
	if det == 0 {
		return nil, false
	}

	det = 1.0 / det
	for i := 0; i < 16; i++ {
		inv[i] = inv[i] * det
	}

	return inv, true
}
