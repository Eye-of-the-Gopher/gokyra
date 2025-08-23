package main

import (
	"fmt"
	"os"

	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
	"github.com/nibrahim/eye-of-the-gopher/pkg/formats"
)

func main() {
	if len(os.Args) != 4 {
		utils.ErrorAndExit("Usage : ./decodecmp dataFile paletteFile ouputImageFile")
		utils.SetupLogging("decode-cmp.log")
	}

	dataFile := os.Args[1]
	paletteFile := os.Args[2]
	ImageFile := os.Args[3]

	CMPData, err := os.ReadFile(dataFile)
	if err != nil {
		utils.ErrorAndExit("Can't read data file %s", dataFile)
	}
	paletteData, err := os.ReadFile(paletteFile)
	if err != nil {
		utils.ErrorAndExit("Can't read palette file %s", paletteFile)
	}
	palette := formats.DecodePalette(paletteData)
	decompressedData := formats.DecodeCmp(dataFile, CMPData, palette)

	fmt.Printf("Writing %s\n", ImageFile)
	utils.WriteCMPToPNG(decompressedData, ImageFile, palette, 320, 200)

}
