package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// Use text as a button
type TextButton struct {
	img         *ebiten.Image
	interactive bool // does it respond to mouse
	*UI
}

// NewButton create a new button with 4 states, normal, hover and press, disable
// to be drawn at posX posY on screen
func NewTextButton(txt string, c color.Color, face font.Face,
	posX, posY int, interactive bool) *TextButton {
	textRect, _ := font.BoundString(face, txt)
	ascent := -textRect.Min.Y.Floor()
	imgW := textRect.Max.X.Ceil() + textRect.Min.X.Ceil()
	imgH := textRect.Max.Y.Ceil() - textRect.Min.Y.Floor()
	img := ebiten.NewImage(imgW, imgH)
	text.Draw(img, txt, face, 0, ascent, c)

	return &TextButton{img: img, interactive: interactive,
		UI: NewUI(posX, posY, imgW, imgH),
	}
}

func (b *TextButton) Draw(screen *ebiten.Image) {
	b.ResetOptions()
	if !b.interactive {
		screen.DrawImage(b.img, b.Op)
	}
	if !b.Active {
		b.Op.ColorScale.ScaleAlpha(0.5)
		screen.DrawImage(b.img, b.Op)
		return
	}
	if b.IsDown() {
		b.Op.ColorScale.Scale(1.2, 1.2, 1.2, 1)
		b.Op.GeoM.Translate(2, 2)
		screen.DrawImage(b.img, b.Op)
		return
	}
	if b.IsMouseOnButton() {
		b.Op.ColorScale.Scale(1.2, 1.2, 1.2, 1)
		b.Op.GeoM.Translate(1, 1)
		screen.DrawImage(b.img, b.Op)
		return
	}
	screen.DrawImage(b.img, b.Op)
}
