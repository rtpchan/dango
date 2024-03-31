// Adopted from ebiten camera example
// https://github.com/hajimehoshi/ebiten/tree/main/examples/camera

package dango

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/math/f64"
)

// Camera projects world to Screen
type Camera struct {
	ViewPort   f64.Vec2 // viewport should be the same as the window size/resolution
	Position   f64.Vec2 // points camera to `Position` in the world
	ZoomFactor int
	Rotation   float64
	matrix     ebiten.GeoM
}

func (c *Camera) String() string {
	return fmt.Sprintf(
		"T: %.1f, R: %.02f, S: %d",
		c.Position, c.Rotation, c.ZoomFactor,
	)
}

// SetViewPort set the size of the window
func (c *Camera) SetViewPort(w, h int) {
	c.ViewPort = [2]float64{float64(w), float64(h)}
}

// SetPosition moves camera to location x, y
func (c *Camera) SetPosition(x, y float64) {
	c.Position = [2]float64{x, y}
}

// GetPosition return x, y
func (c *Camera) GetPosition() (float64, float64) {
	return c.Position[0], c.Position[1]
}

func (c *Camera) Pan(x, y float64) {
	c.Position[0] += x*math.Cos(-c.Rotation) - y*math.Sin(-c.Rotation)
	c.Position[1] += x*math.Sin(-c.Rotation) + y*math.Cos(-c.Rotation)
}

func (c *Camera) Rotate(a float64) {
	c.Rotation += a
}

func (c *Camera) ZoomIn(n int) {
	c.ZoomFactor += n
}

func (c *Camera) ZoomOut(n int) {
	c.ZoomFactor -= n
}

func (c *Camera) viewportCenter() f64.Vec2 {
	return f64.Vec2{
		c.ViewPort[0] * 0.5,
		c.ViewPort[1] * 0.5,
	}
}

// UpdateMatrix when position/rotation/zoom changes
func (c *Camera) Update() {
	c.matrix = c.worldMatrix()
}

// Matrix return current camera matrix
func (c *Camera) Matrix() ebiten.GeoM {
	return c.matrix
}

// Ebiten GeoM,
// usage sprite.GeoM.Concat(camera.GeoM)
func (c *Camera) GeoM() ebiten.GeoM {
	return c.matrix
}

// DEPRECATED,  Concat camera's matrix with m
func (c *Camera) SpriteGeoMConcat(sprite ebiten.GeoM) ebiten.GeoM {
	nm := ebiten.GeoM{}
	nm.Concat(sprite)
	nm.Concat(c.matrix)
	return nm
}

func (c *Camera) worldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.Position[0], -c.Position[1])
	// We want to scale and rotate around center of image / screen
	// m.Translate(-c.Position[0]+c.viewportCenter()[0], -c.Position[1]+c.viewportCenter()[1])
	// // We want to scale and rotate around center of image / screen
	// m.Translate(-c.viewportCenter()[0], -c.viewportCenter()[1])
	m.Scale(
		math.Pow(1.01, float64(c.ZoomFactor)),
		math.Pow(1.01, float64(c.ZoomFactor)),
	)
	m.Rotate(c.Rotation)
	m.Translate(c.viewportCenter()[0], c.viewportCenter()[1])
	return m
}

func (c *Camera) Render(screen, world *ebiten.Image) {
	screen.DrawImage(world, &ebiten.DrawImageOptions{
		GeoM: c.worldMatrix(),
	})
}

func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
	inverseMatrix := c.matrix
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		return inverseMatrix.Apply(float64(posX), float64(posY))
	} else {
		// When scaling it can happend that matrix is not invertable
		return math.NaN(), math.NaN()
	}
}

func (c *Camera) WorldToScreen(wx, wy float64) (float64, float64) {
	sx, sy := c.matrix.Apply(wx, wy)
	return sx, sy
}

// WorldToScreen32 return x, y in float32, ebiten screen drawing use float32
func (c *Camera) WorldToScreen32(wx, wy float64) (float32, float32) {
	sx, sy := c.matrix.Apply(wx, wy)
	return float32(sx), float32(sy)
}

// IsPointInViewport check is a point in on screen
func (c *Camera) IsPointInViewport(wx, wy float64) bool {
	m := c.worldMatrix()
	sx, sy := m.Apply(wx, wy)
	if sx >= 0 && sx <= c.ViewPort[0] && sy >= 0 && sy <= c.ViewPort[1] {
		return true
	}
	return false
}

func (c *Camera) Reset() {
	c.Position[0] = 0
	c.Position[1] = 0
	c.Rotation = 0
	c.ZoomFactor = 0
}
