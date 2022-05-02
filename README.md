# Dango 
**not tested, not production ready, and ABSOLUTELY NO WARRANTY**

Dango is a collection of functions that can be used with [Ebiten](ebiten.org)

## data.go - FS
### embedded file system
Commonly used function to load files from embedded file system
Example:
```
import "github.com/iatearock/dango"
//go:embed assets/*
var data embed.FS
vfs = dango.NewFS(data)

// example usage
img, err := vfs.GetImage("assets/images/red.png") // *ebiten.Image
ff, err := vfs.GetFontFace("assets/font/red.ttf") // font.Face
str, err := vgs.GetCSV("assets/csv/red.csv") // [][]string
```

## camera
Adopted from ebiten camera example:
https://github.com/hajimehoshi/ebiten/tree/main/examples/camera

## scene
Scene manager adopted from ebiten Block example
Handle transition between scenes that implement Update() and Draw(*ebiten.Image)

## vector
Vector maths, **copied** directly from, 
https://github.com/jakecoffman/cp v1.1.0

## id
Simple unique id generator, concurrency safe, I think.