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
