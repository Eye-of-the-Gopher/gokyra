package main

import (
	"log/slog"
	"os"

	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
	"github.com/nibrahim/eye-of-the-gopher/pkg/formats"
	unpak "github.com/nibrahim/eye-of-the-gopher/pkg/formats"
)

func main() {
	formats.InitLogger(formats.AssetLoaderConfig{
		AssetLevel: slog.LevelDebug,
		CmpLevel:   slog.LevelError,
		MazLevel:   slog.LevelError,
		PakLevel:   slog.LevelDebug,
		PalLevel:   slog.LevelError,
	})

	if len(os.Args) != 3 {
		utils.ErrorAndExit("Usage: unpak pakfile outputDirectory")
	}

	assets := unpak.NewAssets()
	err := assets.LoadPakFile(os.Args[1], "")
	if err != nil {
		utils.ErrorAndExit("Could not unpack file: %v", err)
	}
	opdir := os.Args[2]
	assets.WriteAssetData(opdir)

}
