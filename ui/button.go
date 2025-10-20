package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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
	imgSize := img.Bounds().Size()
	imgW := imgSize.X
	imgH := imgSize.Y
	iA := ebiten.NewImage(imgW, imgH)
	iA.DrawImage(img, &ebiten.DrawImageOptions{})
	iH := ebiten.NewImage(imgW, imgH)
	iH.DrawImage(hover, &ebiten.DrawImageOptions{})
	iP := ebiten.NewImage(imgW, imgH)
	iP.DrawImage(press, &ebiten.DrawImageOptions{})
	iD := ebiten.NewImage(imgW, imgH)
	iD.DrawImage(disable, &ebiten.DrawImageOptions{})

	return &Button{img: iA, imgHover: iH, imgPress: iP, imgDisable: iD,
		UI: NewUI(posX, posY, imgW, imgH),
	}
}

// SetText create text on button, text cannot be remove once set.
func (b *Button) SetText(txt string, face text.Face, color color.Color) {
	tw, th := text.Measure(txt, face, 2)
	// textRect, _ := font.BoundString(face, txt)
	// tw := int(textRect.Max.X)
	// th := int(textRect.Max.Y)
	moveX := (float64(b.ImgW) - tw) / 2
	moveY := (float64(b.ImgH)-th)/2 + th
	op := &text.DrawOptions{}
	op.GeoM.Translate(moveX, moveY)
	op.ColorScale.ScaleWithColor(color)
	text.Draw(b.img, txt, face, op)
	text.Draw(b.imgHover, txt, face, op)
	text.Draw(b.imgPress, txt, face, op)
}

// SetText place img on top of the button, img cannot be remove once set.
func (b *Button) SetImage(img *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	b.img.DrawImage(img, op)
	b.imgHover.DrawImage(img, op)
	b.imgPress.DrawImage(img, op)
	b.imgDisable.DrawImage(img, op)
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
