package dango

import (
	"fmt"
	"math"
)

type Camera3D struct {
	pos    Vector3
	lookAt Vector3
	fov    float64 // degree, field of view horizonatlly
	w      float64 // screen size in pixels
	h      float64

	viewMatrix       []float64 // Camera3D matrix
	projectionMatrix []float64 // perspective matrix
	viewportMatrix   []float64
	mvp              []float64 // model-view-perspective
	combineMatrix    []float64 // viewport * projection * view

	changed bool // track if Camera3D requires update
}

func NewCamera3D(pos, lookAt Vector3, fov, w, h float64) *Camera3D {
	cam := &Camera3D{pos: pos, lookAt: lookAt, fov: fov, w: w, h: h, changed: true}
	cam.Update()
	return cam
}

// call after changeing Camera3D parameters
func (cam *Camera3D) Update() {
	if cam.changed {
		cam.UpdateCamera3D(cam.pos, cam.lookAt, cam.fov, cam.w, cam.h)
		cam.changed = false
	}
}

// move the Camera3D
func (cam *Camera3D) Move(dir Vector3) {
	if dir.Length() == 0. {
		return
	}
	moveDir := Vector{cam.lookAt.X - cam.pos.X, cam.lookAt.Z - cam.pos.Z}
	moveForward := moveDir.Normalize().Mult(dir.X)
	moveSideway := moveDir.Perp().Normalize().Mult(dir.Z)
	move := moveForward.Add(moveSideway)

	cam.pos = cam.pos.Add(Vector3{move.X, dir.Y, move.Y})
	cam.lookAt = cam.lookAt.Add(Vector3{move.X, dir.Y, move.Y})

	cam.changed = true
}

// rotate the Camera3D, need to call Update() the Camera3D manually
func (cam *Camera3D) Yaw(rad float64) {
	if rad == 0. {
		return
	}
	vec := cam.lookAt.Sub(cam.pos).Normalize()
	newX := vec.X*math.Cos(rad) - vec.Z*math.Sin(rad)
	newZ := vec.X*math.Sin(rad) + vec.Z*math.Cos(rad)
	cam.lookAt = cam.pos.Add(Vector3{newX, vec.Y, newZ})

	cam.changed = true
}

// pitch the Camera3D, need to call Update() the Camera3D manually
func (cam *Camera3D) Pitch(rad float64) {
	if rad == 0. {
		return
	}
	zAxis := cam.pos.Sub(cam.lookAt).Normalize()
	xAxis := Vector3{0, 1, 0}.Cross(zAxis)
	// Rodrigues' rotation formula
	term1 := zAxis.Mult(math.Cos(rad))
	term2 := xAxis.Cross(zAxis).Mult(math.Sin(rad))
	term3 := xAxis.Mult(xAxis.Dot(zAxis)).Mult(1 - math.Cos(rad))
	v := term1.Add(term2).Add(term3)
	parallel := v.Dot(Vector3{0, 1, 0})
	if math.Abs(parallel) > 0.9 {
		// too close to pointing verticall up or down, do not update
		return
	}
	cam.lookAt = cam.pos.Sub(v)

	cam.changed = true
}

// change fov, need to call Update() the camear manually
func (cam *Camera3D) FOV(delta float64) {
	cam.fov = cam.fov + delta
	cam.changed = true
}

func (cam *Camera3D) UpdateCamera3D(pos, lookAt Vector3, fov, w, h float64) {
	fovX := fov / 180. * math.Pi
	aspectRatio := w / h
	fovY := 2. * math.Atan(math.Tan(fovX/2.)/aspectRatio) // radian
	near := 5.
	far := 500.
	up := Vector3{0, 1, 0}.Normalize()
	// zAxis := pos.Sub(lookAt).Normalize()
	zAxis := lookAt.Sub(pos).Normalize()
	xAxis := up.Cross(zAxis)
	yAxis := zAxis.Cross(xAxis)
	translate := Vector3{X: -xAxis.Dot(pos),
		Y: -yAxis.Dot(pos),
		Z: -zAxis.Dot(pos),
	}
	cam.pos = pos
	// cam.lookAt = cam.pos.Sub(zAxis) // lookAt
	cam.lookAt = lookAt
	cam.fov = fov
	cam.w = w
	cam.h = h

	cam.viewMatrix = []float64{
		xAxis.X, xAxis.Y, xAxis.Z, translate.X,
		yAxis.X, yAxis.Y, yAxis.Z, translate.Y,
		zAxis.X, zAxis.Y, zAxis.Z, translate.Z,
		0, 0, 0, 1,
	}
	cam.projectionMatrix = []float64{
		1. / (aspectRatio * math.Tan(fovY/2.)), 0, 0, 0,
		0, 1. / math.Tan(fovY/2.0), 0, 0,
		// 0, 0, -(far + near) / (far - near), -(2. * far * near) / (far - near),
		// 0, 0, -1, 0,
		0, 0, (far) / (far - near), (-near * far) / (far - near),
		0, 0, 1, 0,
	}
	vpX := 0. // origin of viewport, in screen coordinate
	vpY := 0.
	zNear := 0. // viewing box,
	zFar := 1.  // usually between 0 to 1
	cam.viewportMatrix = []float64{
		w / 2., 0, 0, vpX + w/2.,
		0, -h / 2., 0, vpY + h/2.,
		0, 0, (zFar - zNear) / 2., zNear + (zFar-zNear)/2.,
		0, 0, 0, 1.,
	}
	cam.updateCombineMatrix()
}

func (cam *Camera3D) updateCombineMatrix() {
	// ScreenPos = ViewportMatrix * ProjectionMatrix * ViewMatrix * ModelMatrix * WorldPos
	// projection * view
	p := cam.projectionMatrix
	v := cam.viewMatrix
	cam.mvp = []float64{
		p[0] * v[0], p[0] * v[1], p[0] * v[2], p[0] * v[3],
		p[5] * v[4], p[5] * v[5], p[5] * v[6], p[5] * v[7],
		p[10] * v[8], p[10] * v[9], p[10] * v[10], p[10]*v[11] + p[11]*v[15],
		p[14] * v[8], p[14] * v[9], p[14] * v[10], p[14] * v[11],
	}

	// combineMatrix = viewport * temp (mvp matrix)
	cam.viewportMultMVP()
}

// convert world position to screen coordinate
func (cam *Camera3D) PosToScreen(p Vector3) []float32 {
	c := cam.combineMatrix
	w := c[12]*p.X + c[13]*p.Y + c[14]*p.Z + c[15]
	// divide by w for perspective division, when w is -ve, the point is
	// behind the Camera3D
	x := (c[0]*p.X + c[1]*p.Y + c[2]*p.Z + c[3]) / w
	y := (c[4]*p.X + c[5]*p.Y + c[6]*p.Z + c[7]) / w
	// z, typically between 0 and 1, encodes how deep this point is in the viewport
	z := (c[8]*p.X + c[9]*p.Y + c[10]*p.Z + c[11]) / w

	return []float32{float32(x), float32(y), float32(z), float32(w)}
}

// convert world position to screen coordinate
func (cam *Camera3D) WorldToScreen(p []float64) []float64 {
	c := cam.combineMatrix
	w := c[12]*p[0] + c[13]*p[1] + c[14]*p[2] + c[15]
	// divide by w for perspective division, when w is -ve, the point is
	// behind the Camera3D
	x := (c[0]*p[0] + c[1]*p[1] + c[2]*p[2] + c[3]) / w
	y := (c[4]*p[0] + c[5]*p[1] + c[6]*p[2] + c[7]) / w
	// z, typically between 0 and 1, encodes how deep this point is in the viewport
	z := (c[8]*p[0] + c[9]*p[1] + c[10]*p[2] + c[11]) / w

	return []float64{x, y, z, w}

	// fmt.Println("---world coor---")
	// fmt.Println(p)
	// t := MatrixVectorMultiplication(cam.viewMatrix, p)
	// fmt.Println("---view * world coor-----------")
	// fmt.Println(t)
	// t = MatrixVectorMultiplication(cam.projectionMatrix, t)
	// fmt.Println("--- * projection -----------")
	// fmt.Println(t)
	// t = MatrixVectorMultiplication(cam.viewportMatrix, t)
	// fmt.Println("--- * viewport -----------")
	// fmt.Println(t)
	// t = []float64{t[0] / t[3], t[1] / t[3], t[2] / t[3], t[3] / t[3]}
	// fmt.Println("--- / clip coor -----------")
	// fmt.Println(t)
	// return t
}

// give screen coordinates of 2 points
// {x1 y1 z1 w1 x2 y2 z2 w2}, i.e. point1 [0] v[1], and point2 v[4] v[5]
// if one point is behind the Camera3D, it will return move that point to where
// the line intersect with the clip plan
// if both points are behind, it will return both points  as { 0, 0, 0, -1}
// return boolean indicates if the line is in front of the Camera3D
func (cam *Camera3D) LineToScreen(a Vector3, b Vector3) ([]float32, bool) {
	va := []float64{a.X, a.Y, a.Z, 1.}
	vb := []float64{b.X, b.Y, b.Z, 1.}

	// MVP * point
	pa := MatrixVectorMultiplication(cam.mvp, va)
	pb := MatrixVectorMultiplication(cam.mvp, vb)

	// w < 0, or z < -w
	behindClipA := pa[3] <= 0 || (pa[2] < 0)
	behindClipB := pb[3] <= 0 || (pb[2] < 0)
	if behindClipA && behindClipB {
		// both points behind Camera3D clip plane
		return []float32{float32(pa[0]), float32(pa[1]), float32(pa[2]), float32(pa[3]), float32(pb[0]), float32(pb[1]), float32(pb[2]), float32(pb[3])}, false
	}

	if behindClipA || behindClipB {
		// one of the points is behind Camera3D
		// solve for t where the line pa + t(pb-pa) intersects the z=0 plane
		numerator := pa[2]
		denominator := pa[2] - pb[2]

		if denominator != 0. {
			t := numerator / denominator

			// Calculate the single, unique intersection point on the near plane
			ix := pa[0] + t*(pb[0]-pa[0])
			iy := pa[1] + t*(pb[1]-pa[1])
			iz := pa[2] + t*(pb[2]-pa[2])
			iw := pa[3] + t*(pb[3]-pa[3])

			// Now, replace the vertex that was behind the plane with this new point.
			if behindClipA {
				pa[0] = ix
				pa[1] = iy
				pa[2] = iz
				pa[3] = iw
			} else { // behindClipB must be true
				pb[0] = ix
				pb[1] = iy
				pb[2] = iz
				pb[3] = iw
			}
		}
	}

	// fmt.Println("mvp * pt")
	// fmt.Println(cam.PrintVector(pa))
	// fmt.Println(cam.PrintVector(pb))

	// viewport * point
	vpa := MatrixVectorMultiplication(cam.viewportMatrix, pa)
	vpb := MatrixVectorMultiplication(cam.viewportMatrix, pb)

	// fmt.Println("viewport * pt")
	// fmt.Println(cam.PrintVector(vpa))
	// fmt.Println(cam.PrintVector(vpb))

	// perspective divide
	sa := []float32{float32(vpa[0]) / float32(vpa[3]), float32(vpa[1]) / float32(vpa[3]), float32(vpa[2]) / float32(vpa[3]), float32(vpa[3])}
	sb := []float32{float32(vpb[0]) / float32(vpb[3]), float32(vpb[1]) / float32(vpb[3]), float32(vpb[2]) / float32(vpb[3]), float32(vpb[3])}

	// fmt.Println("output pt")
	// fmt.Println(cam.PrintVector(sa))
	// fmt.Println(cam.PrintVector(sb))

	return []float32{sa[0], sa[1], sa[2], sa[3], sb[0], sb[1], sb[2], sb[3]}, true
}

func (cam *Camera3D) viewportMultMVP() {
	// viewport * temp (mvp matrix)
	v := cam.viewportMatrix
	t := cam.mvp
	cam.combineMatrix = []float64{
		v[0]*t[0] + v[3]*t[12],
		v[0]*t[1] + v[3]*t[13],
		v[0]*t[2] + v[3]*t[14],
		v[0]*t[3] + v[3]*t[15],

		v[5]*t[4] + v[7]*t[12],
		v[5]*t[5] + v[7]*t[13],
		v[5]*t[6] + v[7]*t[14],
		v[5]*t[7] + v[7]*t[15],

		v[10]*t[8] + v[11]*t[12],
		v[10]*t[9] + v[11]*t[13],
		v[10]*t[10] + v[11]*t[14],
		v[10]*t[11] + v[11]*t[15],

		t[12], t[13], t[14], t[15],
	}
}

// print a 4 x 4 matrix
func (cam *Camera3D) PrintMatrix(m []float64) string {
	return fmt.Sprintf("%0.4f, %0.4f, %0.4f, %0.4f\n%0.4f, %0.4f, %0.4f, %0.4f\n%0.4f, %0.4f, %0.4f, %0.4f\n%0.4f, %0.4f, %0.4f, %0.4f\n", m[0], m[1], m[2], m[3], m[4], m[5], m[6], m[7], m[8], m[9], m[10], m[11], m[12], m[13], m[14], m[15])
}

// print a 4 x 1 vector
func (cam *Camera3D) PrintVector(m []float64) string {
	return fmt.Sprintf("%0.4f, %0.4f, %0.4f, %0.4f\n", m[0], m[1], m[2], m[3])
}
func (cam *Camera3D) PrintViewMatrix() string {
	return fmt.Sprintf("%0.4f, %0.4f, %0.4f, %0.4f\n%0.4f, %0.4f, %0.4f, %0.4f\n%0.4f, %0.4f, %0.4f, %0.4f\n%0.4f, %0.4f, %0.4f, %0.4f\n", cam.viewMatrix[0], cam.viewMatrix[1], cam.viewMatrix[2], cam.viewMatrix[3], cam.viewMatrix[4], cam.viewMatrix[5], cam.viewMatrix[6], cam.viewMatrix[7], cam.viewMatrix[8], cam.viewMatrix[9], cam.viewMatrix[10], cam.viewMatrix[11], cam.viewMatrix[12], cam.viewMatrix[13], cam.viewMatrix[14], cam.viewMatrix[15])
}
func (cam *Camera3D) PrintProjectionMatrix() string {
	return fmt.Sprintf("%0.4f, %0.4f, %0.4f, %0.4f\n%0.4f, %0.4f, %0.4f, %0.4f\n%0.4f, %0.4f, %0.4f, %0.4f\n%0.4f, %0.4f, %0.4f, %0.4f\n", cam.projectionMatrix[0], cam.projectionMatrix[1], cam.projectionMatrix[2], cam.projectionMatrix[3], cam.projectionMatrix[4], cam.projectionMatrix[5], cam.projectionMatrix[6], cam.projectionMatrix[7], cam.projectionMatrix[8], cam.projectionMatrix[9], cam.projectionMatrix[10], cam.projectionMatrix[11], cam.projectionMatrix[12], cam.projectionMatrix[13], cam.projectionMatrix[14], cam.projectionMatrix[15])
}
func (cam *Camera3D) PrintViewportMatrix() string {
	return fmt.Sprintf("%0.4f, %0.4f, %0.4f, %0.4f\n%0.4f, %0.4f, %0.4f, %0.4f\n%0.4f, %0.4f, %0.4f, %0.4f\n%0.4f, %0.4f, %0.4f, %0.4f\n", cam.viewportMatrix[0], cam.viewportMatrix[1], cam.viewportMatrix[2], cam.viewportMatrix[3], cam.viewportMatrix[4], cam.viewportMatrix[5], cam.viewportMatrix[6], cam.viewportMatrix[7], cam.viewportMatrix[8], cam.viewportMatrix[9], cam.viewportMatrix[10], cam.viewportMatrix[11], cam.viewportMatrix[12], cam.viewportMatrix[13], cam.viewportMatrix[14], cam.viewportMatrix[15])
}
func (cam *Camera3D) PrintCombineMatrix() string {
	return fmt.Sprintf("%0.4f, %0.4f, %0.4f, %0.4f\n%0.4f, %0.4f, %0.4f, %0.4f\n%0.4f, %0.4f, %0.4f, %0.4f\n%0.4f, %0.4f, %0.4f, %0.4f\n", cam.combineMatrix[0], cam.combineMatrix[1], cam.combineMatrix[2], cam.combineMatrix[3], cam.combineMatrix[4], cam.combineMatrix[5], cam.combineMatrix[6], cam.combineMatrix[7], cam.combineMatrix[8], cam.combineMatrix[9], cam.combineMatrix[10], cam.combineMatrix[11], cam.combineMatrix[12], cam.combineMatrix[13], cam.combineMatrix[14], cam.combineMatrix[15])
}
