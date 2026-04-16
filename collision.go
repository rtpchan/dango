package dango

import "math"

/// Light weight collision detection for rectangles and circles

type Rectangle struct {
	X, Y, W, H float64
}

func (r Rectangle) GetPos() (float64, float64) {
	return r.X, r.Y
}

type Circle struct {
	X, Y, R float64
}

func (c Circle) GetPos() (float64, float64) {
	return c.X, c.Y
}

type ShapeProvider interface {
	GetShape() any // return Rectangle or Circle
}

type ShapePosition interface {
	GetPos() (float64, float64)
}

type Cell struct {
	X, Y int
}

// SpatialGrid stores entity IDs mapped to grid cells
// ID is a user id, any comparable types
// T is a *Entity that has GetShape()
type SpatialGrid[ID comparable, T ShapeProvider] struct {
	CellSize float64
	Data     map[Cell]map[ID]T
}

// NewSpatialGrid initializes a new grid
func NewSpatialGrid[ID comparable, T ShapeProvider](cellSize float64) *SpatialGrid[ID, T] {
	return &SpatialGrid[ID, T]{
		CellSize: cellSize,
		Data:     make(map[Cell]map[ID]T),
	}
}

// posToCell converts world coordinates to grid coordinates
func (g *SpatialGrid[ID, T]) posToCell(x, y float64) Cell {
	return Cell{
		X: int(math.Floor(x / g.CellSize)),
		Y: int(math.Floor(y / g.CellSize)),
	}
}

// Add adds an entity ID to the grid based on its position
func (g *SpatialGrid[ID, T]) Add(id ID, t T) {
	x, y := t.GetShape().(ShapePosition).GetPos()
	cell := g.posToCell(x, y)
	if _, exists := g.Data[cell]; !exists {
		g.Data[cell] = make(map[ID]T)
	}
	g.Data[cell][id] = t
}

// Remove removes an entity ID from the grid
func (g *SpatialGrid[ID, T]) Remove(id ID, t T) {
	x, y := t.GetShape().(ShapePosition).GetPos()
	cell := g.posToCell(x, y)
	if ids, exists := g.Data[cell]; exists {
		delete(ids, id)
		// Clean up the cell if it's now empty
		if len(ids) == 0 {
			delete(g.Data, cell)
		}
	}
}

// Query returns a slice of potential collision shapes in the 9 surrounding cells
func (g *SpatialGrid[ID, T]) Query(x, y float64) []T {
	center := g.posToCell(x, y)
	var results []T

	// Check the 3x3 neighborhood
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			neighborCell := Cell{X: center.X + dx, Y: center.Y + dy}
			if ids, exists := g.Data[neighborCell]; exists {
				for id := range ids {
					results = append(results, g.Data[neighborCell][id])
				}
			}
		}
	}
	return results
}

// Update handles moving an entity from one cell to another
func (g *SpatialGrid[ID, T]) Update(id ID, t T, oldX, oldY float64) {
	x, y := t.GetShape().(ShapePosition).GetPos()
	oldCell := g.posToCell(oldX, oldY)
	newCell := g.posToCell(x, y)

	// Only update if the entity has actually crossed into a new cell
	if oldCell != newCell {
		//Remove from old cell
		if ids, exists := g.Data[oldCell]; exists {
			delete(ids, id)
		}
		g.Add(id, t)
	}
}

// Check if a position collide with anything
func (g *SpatialGrid[ID, T]) Collide(entity T, x, y float64) bool {
	shape := entity.GetShape()
	nearby := g.Query(x, y)
	for _, other := range nearby {
		if Intersects(shape, other.GetShape()) {
			return true
		}
	}
	return false
}

func sq(x float64) float64 {
	return x * x
}

func distSq(x1, y1, x2, y2 float64) float64 {
	return sq(x2-x1) + sq(y2-y1)
}

func clamp(v, minV, maxV float64) float64 {
	return math.Max(minV, math.Min(v, maxV))
}

// --- Intersection Logic ---

// IntersectsRR checks collision between two Rectangles (AABB)
func IntersectsRR(r1, r2 Rectangle) bool {
	return !(r1.X > r2.X+r2.W ||
		r1.X+r1.W < r2.X ||
		r1.Y > r2.Y+r2.H ||
		r1.Y+r1.H < r2.Y)
}

// IntersectsCC checks collision between two Circles
func IntersectsCC(c1, c2 Circle) bool {
	return distSq(c1.X, c1.Y, c2.X, c2.Y) <= sq(c1.R+c2.R)
}

// IntersectsRC checks collision between a Rectangle and a Circle
func IntersectsRC(r Rectangle, c Circle) bool {
	closestX := clamp(c.X, r.X, r.X+r.W)
	closestY := clamp(c.Y, r.Y, r.Y+r.H)

	return distSq(c.X, c.Y, closestX, closestY) <= sq(c.R)
}

// Intersects handles dispatching to the correct collision logic based on types.
// It accepts 'any' to mimic the flexible Clojure API.
func Intersects(a, b any) bool {
	switch s1 := a.(type) {
	case Rectangle:
		switch s2 := b.(type) {
		case Rectangle:
			return IntersectsRR(s1, s2)
		case Circle:
			return IntersectsRC(s1, s2)
		}
	case Circle:
		switch s2 := b.(type) {
		case Rectangle:
			return IntersectsRC(s2, s1) // Order swapped for RC logic
		case Circle:
			return IntersectsCC(s1, s2)
		}
	}
	return false
}
