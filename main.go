package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func setLogger() {
	logger := slog.New(slog.NewTextHandler(
		os.Stderr,
		&slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

}

func main() {
	setLogger()
	slog.Info("Starting program")
	unpakked, err := extractPakFile(os.Args[1])
	var fname string
	if len(os.Args) > 2 {
		fname = os.Args[2]
	}

	if err != nil {
		fmt.Println(err)
	} else {
		for i := range unpakked {
			entry := unpakked[i]
			if filepath.Ext(entry.name) == ".CMP" {
				var data []byte
				if fname != "" && fname != entry.name {
					continue // Skip processing since we're looking for a specific file
				}
				data = decodeCmp(entry.name, entry.data)
				debugFile := fmt.Sprintf("tmp/%s.png", entry.name)
				fmt.Printf("Writing %s\n", debugFile)
				writeCMPToPNG(data, debugFile, 320, 200)
				if fname != "" {
					break // Break out since we've processed the file that we wanted to.
				}
			}
		}

	}

}
