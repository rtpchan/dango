package ui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Digits struct {
	digits []*Digit
}

// NewDigits creates a new number with `num` of digits
func NewDigits(f font.Face, posX, posY int, num int, c color.Color) *Digits {
	ds := []*Digit{}
	rect := text.BoundString(f, "8")
	w := rect.Dx()
	for i := 0; i < num; i++ {
		d := NewDigit(f, posX+int(float64(i*w)*1.3), posY, c)
		ds = append(ds, d)
	}
	return &Digits{digits: ds}
}

func (d *Digits) Update() {
	for _, di := range d.digits {
		di.Update()
	}
}
func getSpecificDigit(num, place int) int {
	r := num % int(math.Pow(10, float64(place)))
	return r / int(math.Pow(10, float64(place-1)))
}

func (d *Digits) SetNumber(n int) {
	count := len(d.digits)
	for i, di := range d.digits {
		sn := getSpecificDigit(n, count-i)
		di.SetNumber(sn)
	}
}

func (d *Digits) GetNumber() int {
	ns := 0
	num := len(d.digits)
	for i, di := range d.digits {
		ns += di.GetNumber() * int(math.Pow(10, float64(num-i-1)))
	}
	return ns
}

func (d *Digits) Draw(screen *ebiten.Image) {
	for _, di := range d.digits {
		di.Draw(screen)
	}
}

type Digit struct {
	number   int // 0 - 9
	fontface font.Face
	c        color.Color

	*UI
}

func NewDigit(f font.Face, posX, posY int, c color.Color) *Digit {
	rect := text.BoundString(f, "8")
	w := rect.Dx()
	h := rect.Dy()
	d := &Digit{number: 0, fontface: f, c: c,
		UI: NewUI(posX, posY, w, h)}
	return d
}

func (d *Digit) SetNumber(x int) {
	d.number = x
}

func (d *Digit) GetNumber() int {
	return d.number
}

func (d *Digit) Draw(screen *ebiten.Image) {
	text.Draw(screen, fmt.Sprintf("%d", d.number), d.fontface, d.PosX, d.PosY+d.ImgH, d.c)
}

func (d *Digit) Update() {
	if d.IsMouseOnButton() {
		_, y := ebiten.Wheel()
		if y >= 1 {
			d.number += 1
			if d.number == 10 {
				d.number = 0
			}
			return
		}
		if y <= -1 {
			d.number -= 1
			if d.number == -1 {
				d.number = 9
			}
			return
		}
	}
}

// TODO scroll to change digit
// TODO click top/bottom half to change digit
// TODO change button to use embedded UI struct
