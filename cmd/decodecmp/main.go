package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log/slog"
	"os"

	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
	"github.com/nibrahim/eye-of-the-gopher/pkg/formats"
)

// Writes extracted CMP bitmap pattern to PNG
func writeCMPToPNG(data []byte, filename string, palette color.Palette, width int, height int) error {
	slog.Debug("Writing PNG", "length", len(data), "to", filename)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill the image with the raw data
	for y := range height {
		for x := range width {
			index := y*width + x
			if index < len(data) {
				img.Set(x, y, palette[data[index]])
			}
		}
	}

	// Create the output file
	file, err := os.Create(filename)
	if err != nil {
		utils.ErrorAndExit("Could not create output file: %v", err)
		return err
	}
	defer file.Close()

	// Encode as PNG
	return png.Encode(file, img)
}

func main() {
	utils.SetupLogging("decode-cmp.log")
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

		writeCMPToPNG(decompressedData, outputFile, palette, 320, 200)
	}

	os.Exit(0)

}
