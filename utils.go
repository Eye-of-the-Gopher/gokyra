package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
)

func writeCMPToPNG(data []byte, filename string, width, height int) error {
	// Create a new grayscale image
	img := image.NewGray(image.Rect(0, 0, width, height))

	// Fill the image with the raw data
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := y*width + x
			if index < len(data) {
				// Use the byte value directly as grayscale intensity
				img.Set(x, y, color.Gray{Y: data[index]})
			}
		}
	}

	// Create the output file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode as PNG
	return png.Encode(file, img)
}

func bytesToBinary(data []byte) string {
	var result strings.Builder
	for i, b := range data {
		if i > 0 {
			result.WriteString(" ") // Space between bytes
		}
		result.WriteString(fmt.Sprintf("%08b", b))
	}
	return result.String()
}
