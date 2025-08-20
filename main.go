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
		&slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

}

func main() {
	setLogger()
	slog.Info("Starting program")
	unpakked, err := extractPakFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
	} else {
		for i := range unpakked {
			entry := unpakked[i]
			if filepath.Ext(entry.name) == ".CMP" {
				decodeCmp(entry.name, entry.data)
			}
		}
		// unpakkedDir := fmt.Sprintf("./%s_files", filepath.Base(os.Args[1]))
		// writePakData(unpakked, unpakkedDir)
		// slog.Info("Written entries to", "dir", unpakkedDir)
	}

}
