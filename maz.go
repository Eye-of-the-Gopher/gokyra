package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func decodeMAZ(mazFileName string, input []byte) []byte {
	ret := []byte{}
	return ret
}

func drawMap(mapData []byte, opFile string) {
	img := image.NewRGBA(image.Rect(0, 0, 500, 500))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	for x := 10; x < 500; x += 10 {
		for y := 10; y < 500; y += 10 {
			img.Set(x, y, color.Black)
		}
	}

	f, _ := os.Create(opFile)
	defer f.Close()
	png.Encode(f, img)

}
