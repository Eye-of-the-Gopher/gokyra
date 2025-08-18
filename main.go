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
	slog.Info("Parsed PAK file")
	if err != nil {
		fmt.Println(err)
	} else {
		unpakkedDir := fmt.Sprintf("./%s_files", filepath.Base(os.Args[1]))
		writePakData(unpakked, unpakkedDir)
		slog.Info("Written entries to", "dir", unpakkedDir)
	}

}
