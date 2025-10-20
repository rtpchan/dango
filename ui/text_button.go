package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Use text as a button
type TextButton struct {
	img         *ebiten.Image
	interactive bool // does it respond to mouse
	*UI
}

// NewButton create a new button with 4 states, normal, hover and press, disable
// to be drawn at posX posY on screen
func NewTextButton(txt string, c color.Color, face text.Face,
	posX, posY int, interactive bool) *TextButton {
	imgW, imgH := text.Measure(txt, face, 2)
	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(c)
	img := ebiten.NewImage(int(imgW), int(imgH))
	text.Draw(img, txt, face, op)

	return &TextButton{img: img, interactive: interactive,
		UI: NewUI(posX, posY, int(imgW), int(imgH)),
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
