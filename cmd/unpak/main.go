package main

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/nibrahim/eye-of-the-gopher/pkg/formats"
)

func setLogger() {
	logger := slog.New(slog.NewTextHandler(
		os.Stderr,
		&slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

}

func main() {
	slog.Info("Starting program")
	unpakked, _ := formats.ExtractPakFile(os.Args[1])
	opdir := os.Args[2]

	for i := range unpakked {
		opfile := path.Join(opdir, unpakked[i].Name)
		fmt.Println("Writing ", opfile)
		os.WriteFile(opfile, unpakked[i].Data, 0644)
	}

}
