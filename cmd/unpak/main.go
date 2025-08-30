package main

import (
	"os"

	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
	unpak "github.com/nibrahim/eye-of-the-gopher/pkg/formats"
)

func main() {
	utils.SetupLogging("debug-unpak.log")
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
