package main

import (
	"log/slog"
	"os"

	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
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
		utils.ErrorAndExit("Usage : ./decodemaz mazFile outputFile")
	}

	dataFile := os.Args[1]
	outputFile := os.Args[2]

	mazData, err := os.ReadFile(dataFile)
	if err != nil {
		utils.ErrorAndExit("Can't read data file %s", dataFile)
	}

	plan := formats.DecodeMAZ(dataFile, mazData)
	// plan.Validate()
	plan.DrawPlan(outputFile)

}
