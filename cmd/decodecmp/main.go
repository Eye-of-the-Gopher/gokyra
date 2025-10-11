package main

import (
	"flag"
	"fmt"
	"image/png"
	"log/slog"
	"os"

	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
	"github.com/nibrahim/eye-of-the-gopher/pkg/formats"
)

func main() {
	formats.InitLogger(formats.AssetLoaderConfig{
		AssetLevel: slog.LevelDebug,
		CmpLevel:   slog.LevelError,
		MazLevel:   slog.LevelError,
		PakLevel:   slog.LevelDebug,
		PalLevel:   slog.LevelError,
	})
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage : %s [options] imageFile1 imageFile2 ...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nArguments:\n")
		fmt.Fprintf(os.Stderr, "  imageFile    List of image files (.CPS, .CMP etc.) to extract\n")
	}

	outputDir := flag.String("outputDir", ".", "Directory to write decompressed files")
	paletteFile := flag.String("paletteFile", "", "Palette file (.PAL, .COL etc.)")
	flag.Parse()
	if *paletteFile == "" {

		flag.Usage()
		utils.ErrorAndExit("Error: Need a paletteFile to use while decompressing")
	}
	if flag.NArg() == 0 {
		utils.ErrorAndExit("Error: No files specified for decompression")
	}

	fmt.Printf("Palette file is %s\n outputDir is %s\n Image files are %v\n", *paletteFile, *outputDir, flag.Args())

	imageFiles := flag.Args()
	paletteData, err := os.ReadFile(*paletteFile)
	if err != nil {
		utils.ErrorAndExit("Can't read palette file %s", *paletteFile)
	}
	palette := formats.DecodePalette(paletteData)
	slog.Debug("Using", "palette", *paletteFile)

	for _, imageFile := range imageFiles {
		outputFile := utils.ImageName(imageFile, "png", *outputDir)
		slog.Debug("Decompressing", "image", imageFile, "to", outputFile)
		CMPData, err := os.ReadFile(imageFile)
		if err != nil {
			utils.ErrorAndExit("Can't read data file %s", imageFile)
		}
		decompressedData := formats.DecodeCmp(imageFile, CMPData, palette)

		img := formats.CMPToImage(decompressedData, palette, 320, 200, 1)

		file, err := os.Create(outputFile)
		if err != nil {
			utils.ErrorAndExit("Could not create output file: %v", err)
		}
		defer file.Close()

		// Encode as PNG
		png.Encode(file, img)
	}

	os.Exit(0)

}

// Create the output file
