package dango

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestIsPointInViewport(t *testing.T) {
	cam := &Camera{}
	cam.ViewPort = [2]float64{800, 400}
	cam.Position = [2]float64{400, 200}

	if cam.IsPointInViewport(20, 20) != true {
		t.Fatalf("Expected point 20,20 in viewport")
	}
	if cam.IsPointInViewport(820, 20) != false {
		t.Fatalf("Expected point 820,20 outside viewport")
	}
}

func TestMatrixConcat(t *testing.T) {
	cam := &Camera{}
	cam.ViewPort = [2]float64{100, 100}
	cam.Position = [2]float64{50, 50}
	cam.Update()

	wx, wy := cam.ScreenToWorld(50, 50)
	if wx != 50 || wy != 50 {
		t.Errorf("Expect (50,50), got (%f, %f)", wx, wy)
	}

	cam.Position = [2]float64{0, 0}
	cam.Update()
	op := ebiten.GeoM{} // sprite GeoM
	op.Translate(10, 5)
	newMatrix := cam.Matrix()
	newMatrix.Concat(op)
	wx2, wy2 := newMatrix.Apply(0, 0)
	if wx2 != 60 || wy2 != 55 {
		t.Errorf("Expect (60,55), got (%f, %f)", wx2, wy2)
	}
}

func TestMatrixApply(t *testing.T) {
	cam := &Camera{}
	cam.ViewPort = [2]float64{100, 100}
	cam.Position = [2]float64{0, 0}
	cam.Update()

	// op := ebiten.GeoM{} // sprite GeoM
	// op.Translate(10, 5)
	newMatrix := cam.Matrix()
	// newMatrix.Concat(op)
	wx2, wy2 := newMatrix.Apply(40, 30)
	if wx2 != 90 || wy2 != 80 {
		t.Errorf("Expect (90,80), got (%f, %f)", wx2, wy2)
	}
}
