package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

// Convert raw palette-indexed data to PNG
func rawToPNG(rawData []byte, palette color.Palette, width, height int, outputFile string) error {
	// Create a paletted image
	img := image.NewPaletted(image.Rect(0, 0, width, height), palette)

	// Fill the image with raw data
	expectedSize := width * height
	if len(rawData) < expectedSize {
		return fmt.Errorf("raw data too small: got %d bytes, need %d", len(rawData), expectedSize)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := y*width + x
			if index < len(rawData) {
				img.SetColorIndex(x, y, rawData[index])
			}
		}
	}

	// Write PNG file
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

func main() {
	if len(os.Args) != 6 {
		fmt.Printf("Usage: %s <raw_file> <palette_file> <width> <height> <output.png>\n", os.Args[0])
		fmt.Println("Example: ./raw2png op.raw EOBPAL.COL 320 200 westwood_output.png")
		os.Exit(1)
	}

	rawFile := os.Args[1]
	paletteFile := os.Args[2]
	width := parseInt(os.Args[3])
	height := parseInt(os.Args[4])
	outputFile := os.Args[5]
	paletteData, _ := os.ReadFile(paletteFile)
	// Read raw data
	rawData, err := os.ReadFile(rawFile)
	if err != nil {
		fmt.Printf("Error reading raw file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Read %d bytes of raw data\n", len(rawData))

	// Read palette
	palette := decodePalette(paletteData)

	fmt.Printf("Read palette with %d colors\n", len(palette))

	// Convert to PNG
	err = rawToPNG(rawData, palette, width, height, outputFile)
	if err != nil {
		fmt.Printf("Error creating PNG: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created %s (%dx%d)\n", outputFile, width, height)
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

// Compile: go build -o raw2png raw_to_png.go
// Usage: ./raw2png op.raw my_assets/4/EOBPAL.COL 320 200 westwood_output.png
