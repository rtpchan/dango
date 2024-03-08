package dango

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type EmbedFS interface {
	fs.ReadDirFS
	fs.ReadFileFS
}

type FS struct {
	filesystem EmbedFS
}

func NewFS(filesystem EmbedFS) *FS {
	return &FS{filesystem: filesystem}
}

func (f *FS) GetImage(path string) (*ebiten.Image, error) {
	imgByte, err := f.filesystem.ReadFile(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

func (f *FS) MustGetImage(path string) *ebiten.Image {
	img, err := f.GetImage(path)
	if err != nil {
		panic(fmt.Sprintf("Cannot find %s", path))
	}
	return img
}

func (f *FS) GetRGBA(path string) (*image.RGBA, error) {
	imgByte, err := f.filesystem.ReadFile(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		return nil, err
	}
	if dst, ok := img.(*image.RGBA); ok {
		return dst, nil
	}
	b := img.Bounds()
	dst := image.NewRGBA(b)
	draw.Draw(dst, dst.Bounds(), img, b.Min, draw.Src)
	return dst, nil
}

func (f *FS) GetFontFace(path string, size, dpi float64) (font.Face, error) {
	var err error
	fontfile, err := f.filesystem.ReadFile(path)
	if err != nil {
		return nil, err
	}
	tt, err := opentype.Parse(fontfile)
	if err != nil {
		return nil, err
	}
	fontface, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size: size, DPI: dpi, Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}
	return fontface, nil
}

func (f *FS) MustGetFontFace(path string, size, dpi float64) font.Face {
	ff, err := f.GetFontFace(path, size, dpi)
	if err != nil {
		panic(fmt.Sprintf("Cannot find %s", path))
	}
	return ff
}

// ReadCSV return list of rows, only use for small file
func (f *FS) ReadCSV(path string) ([][]string, error) {
	b, err := f.filesystem.ReadFile(path)
	if err != nil {
		log.Println(err)
	}
	r := csv.NewReader(bytes.NewReader(b))
	return r.ReadAll()
}

func (f *FS) MustReadCSV(path string) [][]string {
	ff, err := f.ReadCSV(path)
	if err != nil {
		panic(fmt.Sprintf("Cannot find %s", path))
	}
	return ff
}

func (f *FS) Open(path string) (fs.File, error) {
	return f.filesystem.Open(path)
}
func (f *FS) ReadDir(path string) ([]fs.DirEntry, error) {
	return f.filesystem.ReadDir(path)
}
func (f *FS) ReadFile(path string) ([]byte, error) {
	return f.filesystem.ReadFile(path)
}

func (f *FS) MustReadFile(path string) []byte {
	ff, err := f.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Cannot find %s", path))
	}
	return ff
}
