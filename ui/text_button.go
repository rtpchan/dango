package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type TextButton struct {
	img *ebiten.Image
	*UI
}

// NewButton create a new button with 4 states, normal, hover and press, disable
// to be drawn at posX posY on screen
func NewTextButton(txt string, c color.Color, face font.Face,
	posX, posY int, centre bool) *TextButton {

	textRect, _ := font.BoundString(face, txt)
	imgW := textRect.Max.X.Ceil() - textRect.Min.X.Floor()
	imgH := textRect.Max.Y.Ceil() - textRect.Min.Y.Floor()
	// log.Println(textRect)
	// log.Println(textadv)
	// log.Println(textRect.Max.X)
	// log.Println(textRect.Max.X.Ceil())
	img := ebiten.NewImage(imgW, imgH)
	// img.Fill(color.RGBA{2, 2, 2, 255})
	leftBearing := textRect.Min.X.Ceil()
	ascent := -textRect.Min.Y.Ceil()
	text.Draw(img, txt, face, leftBearing, ascent, c)

	return &TextButton{img: img,
		UI: NewUI(posX, posY, imgW, imgH),
	}
}

func (b *TextButton) Draw(screen *ebiten.Image) {
	b.ResetOptions()
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
