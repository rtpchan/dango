package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Drawable interface {
	Draw(*ebiten.Image)
	SetPos(int, int)
	GetPos() (int, int)
	GetSize() (int, int)
}

// Common UI variables and functions
type UI struct {
	Active bool
	PosX   int
	PosY   int
	ImgW   int
	ImgH   int
	Op     *ebiten.DrawImageOptions
}

func NewUI(posX, posY, imgW, imgH int) *UI {
	gui := UI{PosX: posX, PosY: posY, ImgW: imgW, ImgH: imgH, Active: true}
	option := &ebiten.DrawImageOptions{}
	option.GeoM.Translate(float64(posX), float64(posY))
	gui.Op = option
	return &gui
}

// SetPos put button in new position x, y
func (b *UI) SetPos(x, y int) {
	b.PosX = x
	b.PosY = y
	option := &ebiten.DrawImageOptions{}
	option.GeoM.Translate(float64(x), float64(y))
	b.Op = option
}

func (b *UI) IsMouseOnButton() bool {
	mx, my := ebiten.CursorPosition()
	if mx >= b.PosX && mx <= (b.PosX+b.ImgW) &&
		my >= b.PosY && my <= (b.PosY+b.ImgH) {
		return true
	}
	return false
}

func (b *UI) IsHover() bool {
	return b.IsMouseOnButton()
}

func (b *UI) IsJustReleased() bool {
	if b.Active &&
		inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) &&
		b.IsMouseOnButton() {
		return true
	}
	return false
}

func (b *UI) IsJustPressed() bool {
	if b.Active &&
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) &&
		b.IsMouseOnButton() {
		return true
	}
	return false
}

func (b *UI) IsDown() bool {
	if b.Active &&
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) &&
		b.IsMouseOnButton() {
		return true
	}
	return false
}

func (b *UI) SetActive(x bool) {
	b.Active = x
}

func (b *UI) ResetOptions() {
	option := &ebiten.DrawImageOptions{}
	option.GeoM.Translate(float64(b.PosX), float64(b.PosY))
	b.Op = option
}

func (b *UI) GetPos() (int, int) {
	return b.PosX, b.PosY
}

func (b *UI) GetSize() (int, int) {
	return b.ImgW, b.ImgH
}
