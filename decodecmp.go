package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

func main() {
	logFile, err := os.OpenFile("op.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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

	data, _ := os.ReadFile(os.Args[1])
	decompressedData := decodeCmp(os.Args[1], data)
	debugFile := os.Args[2]
	fmt.Printf("Writing %s\n", debugFile)
	writeCMPToPNG(decompressedData, debugFile, 320, 200)

}
