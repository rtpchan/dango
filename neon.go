package dango

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
)

// Neon add neon light effect of an image,
// Neon adds color `c` within `light` pixels away from any pixels
// with 255 for alpha channel.
// Gaussian blur is applied with sqaure of width `blur`, and `sigma`
// as Gaussian standard deviation
// If `origin`, original image is draw on top of new image
// If `resize`, new image will be larger due to Gaussian effect
// on the edge
func Neon(img *image.RGBA, light, blur int, sigma float64,
	c color.RGBA, origin, resize bool) *image.RGBA {

	imgNeon := AddWidth(img, light, c)
	imgBlur := Blur(imgNeon, blur, Kernal(blur, sigma))
	if !resize {
		imgW := imgBlur.Bounds().Dx()
		edge := (imgW - img.Bounds().Dx()) / 2
		imgTmp := image.NewRGBA(img.Bounds())
		draw.Draw(imgTmp, img.Bounds(), imgBlur,
			image.Point{edge, edge}, draw.Over)
		imgBlur = nil
		imgBlur = imgTmp
	}

	if origin {
		imgW, imgH := imgBlur.Bounds().Dx(), imgBlur.Bounds().Dy()
		edge := (imgW - img.Bounds().Dx()) / 2
		rect := image.Rect(edge, edge, imgW-edge, imgH-edge)
		draw.Draw(imgBlur, rect, img, image.Point{}, draw.Over)
		return imgBlur
	}
	return imgBlur
}

// G calculates kernal value with Gaussian distribution at point x, y
func G(x, y, sigma float64) float64 {
	a := 1. / (2. * math.Pi * sigma * sigma)
	b := math.Exp(-(x*x + y*y) / (2. * sigma * sigma))
	return a * b
}

// Create a gaussian kernal, each row is appended to previous row
func Kernal(size int, sigma float64) []float64 {
	g := []float64{}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			g = append(g, G(float64(i-(size/2)), float64(j-(size/2)), sigma))
		}
	}
	sum := 0.
	for _, gg := range g {
		sum += gg
	}
	normalised := []float64{}
	for _, gg := range g {
		normalised = append(normalised, gg/sum)
	}
	return normalised
}

// Return the RGBA values of kernal with `size` x `size` in a list
func Convolution(img *image.RGBA, x, y, size int) []color.RGBA {
	z := size / 2
	rgba := []color.RGBA{}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			c := img.RGBAAt(x+i-z, y+j-z)
			rgba = append(rgba, c)
		}
	}
	return rgba
}

// Return a new image applying Gaussian Blur
// Note, for my purpose, the reutrn image size is larger
// by the size of the kernal, to capture the entire bluring effect
func Blur(img *image.RGBA, size int, kernal []float64) *image.RGBA {

	newImg := image.NewRGBA(img.Bounds())
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			conv := Convolution(img, i, j, size)
			sum := color.RGBA{0, 0, 0, 0}
			for idx, c := range conv {
				sum.R += uint8(float64(c.R) * kernal[idx])
				sum.G += uint8(float64(c.G) * kernal[idx])
				sum.B += uint8(float64(c.B) * kernal[idx])
				sum.A += uint8(float64(c.A) * kernal[idx])
			}
			newImg.SetRGBA(i, j, sum)
		}
	}
	return newImg
}

// DrawRect with width, height, thickness, tansparent edge
func DrawRect(w, h, t, edge int, c color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// top
	for i := edge; i < w-edge; i++ {
		for k := edge; k-edge < t; k++ {
			img.SetRGBA(i, k, c)
		}
	}
	// bottom
	for i := edge; i < w-edge; i++ {
		for k := edge; k-edge < t; k++ {
			img.SetRGBA(i, h-k-1, c)
		}
	}
	// left
	for j := edge; j < h-edge; j++ {
		for k := edge; k-edge < t; k++ {
			img.SetRGBA(k, j, c)
		}
	}
	// right
	for j := edge; j < h-edge; j++ {
		for k := edge; k-edge < t; k++ {
			img.SetRGBA(w-k-1, j, c)
		}
	}
	return img
}

// AddWidth light effect to any pixel with 255 in alpha channle
// within `dist` from the pixels
func AddWidth(img *image.RGBA, dist int, c color.RGBA) *image.RGBA {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	newImg := image.NewRGBA(image.Rect(0, 0, w+dist*2, h+dist*2))

	tmpC := c
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			if img.RGBAAt(i, j).A > 25 {
				tmpC.A = uint8(math.Min(float64(c.A), float64(img.RGBAAt(i, j).A)))
				ni, nj := i+dist, j+dist // new image coordinates
				Paint(newImg, ni, nj, dist, tmpC)
			}
		}
	}
	return newImg
}

// Paint add colour to pixel radius away from x, y
func Paint(img *image.RGBA, x, y, radius int, c color.RGBA) {
	z := radius*2 + 1
	for i := 0; i < z; i++ {
		for j := 0; j < z; j++ {
			m := x - radius + i
			n := y - radius + j
			d := math.Sqrt(math.Pow(float64(m-x), 2.) +
				math.Pow(float64(n-y), 2.))
			exceed := d - float64(radius)
			imgC := img.RGBAAt(m, n)
			if exceed <= 0. {
				// tmpC := color.RGBA{c.R, c.G, c.B, c.A}
				// img.SetRGBA(m, n, c)
				o := SimpleAlphaComposite(c, imgC)
				img.SetRGBA(m, n, o)
			} else if exceed < 1. {
				tmpC := color.RGBA{c.R, c.G, c.B, uint8(exceed * 255)}
				o := SimpleAlphaComposite(tmpC, imgC)
				img.SetRGBA(m, n, o)
			}
		}
	}
}

// Simple Alpha Composite source with background pixel
// https://www.w3.org/TR/compositing-1/#simplealphacompositing
func SimpleAlphaComposite(s, b color.RGBA) color.RGBA {
	var o color.RGBA
	o.R = tI(tF(s.R)*tF(s.A) + tF(b.R)*tF(b.A)*(1.-tF(s.A)))
	o.G = tI(tF(s.G)*tF(s.A) + tF(b.G)*tF(b.A)*(1.-tF(s.A)))
	o.B = tI(tF(s.B)*tF(s.A) + tF(b.B)*tF(b.A)*(1.-tF(s.A)))
	o.A = tI(tF(s.A) + tF(b.A)*(1.-tF(s.A)))
	return o
}

// from uint8 (0-255) to float64 0.- 1.
func tF(i uint8) float64 {
	return float64(i) / 255.
}

// from float64 0. - 1. to uint8 (0-255)
func tI(f float64) uint8 {
	if f > 1. {
		return 255
	}
	if f < 0. {
		return 0
	}
	return uint8(f * 255.)
}

// Write img to path
func WritePNG(path string, img image.Image) error {
	// save file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

// Read img from path
func LoadPNG(name string) *image.RGBA {
	reader, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	tmp, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	b := tmp.Bounds()
	img := image.NewRGBA(b)
	draw.Draw(img, img.Bounds(), tmp, b.Min, draw.Src)

	return img
}
