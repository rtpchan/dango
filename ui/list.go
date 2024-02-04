package ui

const (
	Horizontal = iota
	Vertical
)

type List struct {
	uis       []*UI
	direction int
}

func NewList(dir int) *List {
	return &List{uis: []*UI{}, direction: dir}
}

func (s *List) AddUI(ui *UI) {
	s.uis = append(s.uis, ui)
}
