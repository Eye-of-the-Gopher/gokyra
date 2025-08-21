package main

import (
	"fmt"
	"log/slog"
	"os"
	"path"
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
	unpakked, _ := extractPakFile(os.Args[1])
	opdir := os.Args[2]

	for i := range unpakked {
		opfile := path.Join(opdir, unpakked[i].name)
		fmt.Println("Writing ", opfile)
		os.WriteFile(opfile, unpakked[i].data, 0644)
	}

}
