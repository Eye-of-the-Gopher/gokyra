package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

func setupLogging() {
	logFile, err := os.OpenFile("op.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		errorAndExit("Couldnt open log file")
		// handle error
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

func main() {
	fmt.Println("./decodecmp dataFile paletteFile ouputImageFile")
	if len(os.Args) != 4 {
		errorAndExit("Usage : ./decodecmp dataFile paletteFile ouputImageFile")
	setupLogging()
	}

	dataFile := os.Args[1]
	paletteFile := os.Args[2]
	ImageFile := os.Args[3]

	CMPData, err := os.ReadFile(dataFile)
	if err != nil {
		errorAndExit("Can't read data file %s", dataFile)
	}
	paletteData, err := os.ReadFile(paletteFile)
	if err != nil {
		errorAndExit("Can't read palette file %s", paletteFile)
	}
	palette := decodePalette(paletteData)
	decompressedData := decodeCmp(dataFile, CMPData, palette)

	fmt.Printf("Writing %s\n", ImageFile)
	writeCMPToPNG(decompressedData, ImageFile, palette, 320, 200)

}
