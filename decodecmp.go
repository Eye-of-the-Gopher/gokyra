package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

func main() {
	logFile, err := os.OpenFile("op.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		// handle error
	}
	defer logFile.Close()

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
	fmt.Println("./decodecmp dataFile paletteFile ouputImageFile")

	dataFile := os.Args[1]
	paletteFile := os.Args[2]
	ImageFile := os.Args[3]

	CMPData, _ := os.ReadFile(dataFile)
	paletteData, _ := os.ReadFile(paletteFile)
	palette := decodePalette(paletteData)
	decompressedData := decodeCmp(dataFile, CMPData, palette)

	fmt.Println(palette)

	fmt.Printf("Writing %s\n", ImageFile)
	writeCMPToPNG(decompressedData, ImageFile, palette, 320, 200)

}
