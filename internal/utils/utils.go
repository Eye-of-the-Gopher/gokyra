package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log/slog"
	"os"
	"strings"
)

func SetupLogging(logfile string) {
	logFile, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		ErrorAndExit("Couldnt open log file: %v", err)
	}

	multiWriter := io.MultiWriter(os.Stderr, logFile)

	logger := slog.New(slog.NewTextHandler(
		multiWriter,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.Attr{}
				}
				return a
			},
		}))
	slog.SetDefault(logger)
}

func WriteCMPToPNG(data []byte, filename string, palette color.Palette, width int, height int) error {
	// Create a new grayscale image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill the image with the raw data
	for y := range height {
		for x := range width {
			index := y*width + x
			if index < len(data) {
				// Use the byte value directly as grayscale intensity
				img.Set(x, y, palette[data[index]])
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

func BytesToBinary(data []byte) string {
	var result strings.Builder
	for i, b := range data {
		if i > 0 {
			result.WriteString(" ") // Space between bytes
		}
		result.WriteString(fmt.Sprintf("%08b", b))
	}
	return result.String()
}

func ErrorAndExit(message string, args ...any) {
	fmt.Fprintf(os.Stderr, message, args...)
	os.Exit(-1)
}
