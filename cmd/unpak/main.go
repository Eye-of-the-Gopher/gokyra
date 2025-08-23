package main

import (
	"fmt"
	"os"
	"path"

	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
	"github.com/nibrahim/eye-of-the-gopher/pkg/formats"
)

func main() {
	utils.SetupLogging("debug-unpak.log")
	if len(os.Args) == 3 {
		utils.ErrorAndExit("Usage: unpak pakfile outputDirectory")
	}

	unpakked, err := formats.ExtractPakFile(os.Args[1])
	if err != nil {
		utils.ErrorAndExit("Could not unpack file: %v", err)
	}
	opdir := os.Args[2]

	for i := range unpakked {
		opfile := path.Join(opdir, unpakked[i].Name)
		fmt.Println("Writing ", opfile)
		os.WriteFile(opfile, unpakked[i].Data, 0644)
	}

}
