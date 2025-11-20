package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log/slog"
	"os"

	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func debugPalette(palette color.Palette, filename string) error {
	const (
		columns      = 8
		rows         = 32
		swatchWidth  = 80
		swatchHeight = 20
		textWidth    = 30
		totalWidth   = columns * (swatchWidth + textWidth)
		totalHeight  = rows * swatchHeight
	)

	// Create image
	img := image.NewRGBA(image.Rect(0, 0, totalWidth, totalHeight))

	// Fill with white background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)

	// Draw each color swatch with index
	for i, col := range palette {
		if i >= 256 {
			break
		}

		column := i / rows
		row := i % rows

		x := column * (swatchWidth + textWidth)
		y := row * swatchHeight

		// Draw color swatch
		swatchRect := image.Rect(x+textWidth, y, x+textWidth+swatchWidth, y+swatchHeight)
		draw.Draw(img, swatchRect, &image.Uniform{col}, image.Point{}, draw.Src)

		// Draw index text
		drawText(img, fmt.Sprintf("%3d", i), x+2, y+swatchHeight/2+3, color.RGBA{0, 0, 0, 255})
	}

	// Save to file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer file.Close()

	return png.Encode(file, img)
}

func drawText(img *image.RGBA, text string, x, y int, col color.RGBA) {
	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)
}

func main() {
	formats.InitLogger(formats.AssetLoaderConfig{
		AssetLevel: slog.LevelDebug,
		CmpLevel:   slog.LevelError,
		MazLevel:   slog.LevelError,
		PakLevel:   slog.LevelDebug,
		PalLevel:   slog.LevelError,
	})
	if len(os.Args) != 3 {
		utils.ErrorAndExit("Usage: ./showpalette dataFile paletteFile ouputImageFile")
	}

	paletteFile := os.Args[1]
	opFile := os.Args[2]

	paletteData, _ := os.ReadFile(paletteFile)
	palette := formats.DecodePalette(paletteData)

	debugPalette(palette, opFile)

}
