package ui

import "github.com/hajimehoshi/ebiten/v2"

type UIDirection int

const (
	Horizontal UIDirection = iota
	Vertical
)

// Draw a list of UI
type List struct {
	uis       []Drawable
	PosX      int
	PosY      int
	direction UIDirection
	spacing   int // pixel between each UI
}

func NewList(x, y int, dir UIDirection) *List {
	return &List{uis: []Drawable{}, direction: dir,
		PosX: x, PosY: y, spacing: 10}
}

func (s *List) AddUI(ui Drawable) {
	count := len(s.uis)
	s.uis = append(s.uis, ui)
	if s.direction == Vertical {
		if count == 0 {
			ui.SetPos(s.PosX, s.PosY)
		} else {
			_, previousY := s.uis[count-1].GetPos()
			_, previousH := s.uis[count-1].GetSize()
			ui.SetPos(s.PosX, previousY+previousH+s.spacing)
		}
	} else if s.direction == Horizontal {
		if count == 0 {
			ui.SetPos(s.PosX, s.PosY)
		} else {
			previousX, _ := s.uis[count-1].GetPos()
			previousW, _ := s.uis[count-1].GetSize()
			ui.SetPos(previousX+previousW+s.spacing, s.PosY)
		}
	}

}

func (s *List) Draw(screen *ebiten.Image) {
	for _, ui := range s.uis {
		ui.Draw(screen)
	}
}

// return index of hover UI, return -1 when not hovering
func (s *List) IsHover() int {
	for i, ui := range s.uis {
		if ui.IsHover() {
			return i
		}
	}
	return -1
}

// return index of mouse down UI, -1 if none
func (s *List) IsDown() int {
	for i, ui := range s.uis {
		if ui.IsDown() {
			return i
		}
	}
	return -1
}

// return index of mouse just pressed UI, -1 if none
func (s *List) IsJustPressed() int {
	for i, ui := range s.uis {
		if ui.IsJustPressed() {
			return i
		}
	}
	return -1
}

// return index of mouse just released UI, -1 if none
func (s *List) IsJustReleased() int {
	for i, ui := range s.uis {
		if ui.IsJustReleased() {
			return i
		}
	}
	return -1
}
