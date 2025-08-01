package dango

import "math"

type Vector3 struct {
	X, Y, Z float64
}

func (a Vector3) ToSlice() []float64 {
	return []float64{a.X, a.Y, a.Z, 1}
}

// take first 3 element as x, y, z
func SliceToVector(a []float64) Vector3 {
	return Vector3{X: a[0], Y: a[1], Z: a[2]}
}

// a add b
func (a Vector3) Add(b Vector3) Vector3 {
	return Vector3{X: (a.X + b.X), Y: (a.Y + b.Y), Z: (a.Z + b.Z)}
}

// a subtract b
func (a Vector3) Sub(b Vector3) Vector3 {
	return Vector3{X: (a.X - b.X), Y: (a.Y - b.Y), Z: (a.Z - b.Z)}
}

// a * float64
func (a Vector3) Mult(t float64) Vector3 {
	return Vector3{X: a.X * t, Y: a.Y * t, Z: a.Z * t}
}

func (a Vector3) Dot(b Vector3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// a cross b
func (a Vector3) Cross(b Vector3) Vector3 {
	return Vector3{X: (a.Y*b.Z - a.Z*b.Y),
		Y: (a.Z*b.X - a.X*b.Z),
		Z: (a.X*b.Y - a.Y*b.X)}
}

func (v Vector3) Negate() Vector3 {
	return Vector3{X: -v.X, Y: -v.Y, Z: -v.Z}
}

func (v Vector3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector3) LengthSq() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vector3) Normalize() Vector3 {
	length := v.Length()
	if length == 0. {
		return v
	}
	return Vector3{X: v.X / length, Y: v.Y / length, Z: v.Z / length}
}

func (a Vector3) Lerp(b Vector3, t float64) Vector3 {
	return a.Add(b.Sub(a).Mult(t))
}

func (a Vector3) DistanceSq(b Vector3) float64 {
	return (a.X-b.X)*(a.X-b.X) + (a.Y-b.Y)*(a.Y-b.Y) + (a.Z-b.Z)*(a.Z-b.Z)
}

// Slerp performs Spherical Linear Interpolation between two vectors.
// It's ideal for smoothly rotating a direction vector.
// Both vectors should be normalized for the best results.
func (a Vector3) Slerp(b Vector3, t float64) Vector3 {
	// Clamp t to the range [0, 1]
	if t < 0.0 {
		t = 0.0
	} else if t > 1.0 {
		t = 1.0
	}

	dot := a.Dot(b)

	// Clamp dot to the range [-1, 1] to handle potential floating-point inaccuracies
	if dot > 1.0 {
		dot = 1.0
	} else if dot < -1.0 {
		dot = -1.0
	}

	// Angle between the two vectors
	theta := math.Acos(dot)

	// If the angle is very small, use linear interpolation to avoid division by zero
	if theta < 1e-6 {
		return a.Lerp(b, t)
	}

	sinTheta := math.Sin(theta)

	// Calculate the scale factors for the interpolation
	scaleA := math.Sin((1.0-t)*theta) / sinTheta
	scaleB := math.Sin(t*theta) / sinTheta

	// Return the interpolated vector
	return a.Mult(scaleA).Add(b.Mult(scaleB))
}

// Angle returns the angle in radians between two vectors.
func (a Vector3) Angle(b Vector3) float64 {
	// Ensure vectors are normalized to get the cosine of the angle directly from the dot product.
	normA := a.Normalize()
	normB := b.Normalize()

	dot := normA.Dot(normB)

	// Clamp dot to the range [-1, 1] to handle potential floating-point inaccuracies
	if dot > 1.0 {
		dot = 1.0
	} else if dot < -1.0 {
		dot = -1.0
	}

	return math.Acos(dot)
}
