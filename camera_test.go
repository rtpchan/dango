package dango

import "testing"

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
