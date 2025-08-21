package main

import (
	"fmt"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(
		os.Stderr,
		&slog.HandlerOptions{Level: slog.LevelDebug,
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
