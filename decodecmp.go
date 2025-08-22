package main

import (
	"fmt"
	"os"
)

func main() {
	setupLogging("cmp-decode.log")
	fmt.Println("./decodecmp dataFile paletteFile ouputImageFile")
	if len(os.Args) != 4 {
		errorAndExit("Usage : ./decodecmp dataFile paletteFile ouputImageFile")
	}

	dataFile := os.Args[1]
	paletteFile := os.Args[2]
	ImageFile := os.Args[3]

	CMPData, err := os.ReadFile(dataFile)
	if err != nil {
		errorAndExit("Can't read data file %s", dataFile)
	}
	paletteData, err := os.ReadFile(paletteFile)
	if err != nil {
		errorAndExit("Can't read palette file %s", paletteFile)
	}
	palette := decodePalette(paletteData)
	decompressedData := decodeCmp(dataFile, CMPData, palette)

	fmt.Printf("Writing %s\n", ImageFile)
	writeCMPToPNG(decompressedData, ImageFile, palette, 320, 200)

}
