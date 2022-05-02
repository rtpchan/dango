package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Button struct {
	img        *ebiten.Image
	imgHover   *ebiten.Image
	imgPress   *ebiten.Image
	imgDisable *ebiten.Image

	*UI
}

// NewButton create a new button with 4 images of 4 states, normal, hover and press, disable
// to be drawn at posX posY on screen
func NewButton(img, hover, press, disable *ebiten.Image, posX, posY int) *Button {
	w, h := img.Size()
	option := &ebiten.DrawImageOptions{}
	option.GeoM.Translate(float64(posX), float64(posY))
	return &Button{img: img, imgHover: hover, imgPress: press, imgDisable: disable,
		UI: NewUI(posX, posY, w, h),
	}
}

// SetText create text on button, text cannot be remove once set.
func (b *Button) SetText(txt string, face font.Face, color color.Color) {
	textRect := text.BoundString(face, txt)
	tw := textRect.Dx()
	th := textRect.Dy()
	moveX := (b.ImgW - tw) / 2
	moveY := (b.ImgH-th)/2 + th
	text.Draw(b.img, txt, face, moveX, moveY, color)
	text.Draw(b.imgHover, txt, face, moveX, moveY, color)
	text.Draw(b.imgPress, txt, face, moveX, moveY, color)
}

func (b *Button) Draw(screen *ebiten.Image) {
	if !b.Active {
		screen.DrawImage(b.imgDisable, b.Op)
		return
	}
	if b.IsDown() {
		screen.DrawImage(b.imgPress, b.Op)
		return
	}
	if b.IsMouseOnButton() {
		screen.DrawImage(b.imgHover, b.Op)
		return
	}
	screen.DrawImage(b.img, b.Op)
}
