package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nibrahim/eye-of-the-gopher/internal/engine"
	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
)

func main() {
	formats.InitLogger(formats.AssetLoaderConfig{
		AssetLevel: slog.LevelError,
		CmpLevel:   slog.LevelError,
		MazLevel:   slog.LevelError,
		PakLevel:   slog.LevelError,
		PalLevel:   slog.LevelError,
	})

	engine.InitLogger(slog.LevelDebug)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage : %s [options] assetDirectory ...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nArguments:\n")
		fmt.Fprintf(os.Stderr, "  assetDirectory    Directory with original EOB .pak files (6 of them)\n")

	}
	scale := flag.Float64("scale", 4.0, "Scaling for all assets")
	extraAssetDir := flag.String("extraAssetDir", "", "Directory to side load extra assets")
	enhanced := flag.Bool("enhanced", false, "Use Side loaded enhanced assets")
	flag.Parse()

	if flag.NArg() == 0 {
		utils.ErrorAndExit("Error: No EOB origin asset directory specified")
	}

	assetDir := flag.Args()[0]
	game := engine.NewGame(assetDir, *extraAssetDir, *enhanced)

	sw := int(float64(engine.ScreenWidth) * *scale)
	sh := int(float64(engine.ScreenHeight) * *scale)
	ebiten.SetWindowSize(sw, sh)
	ebiten.SetWindowTitle("Eye Of The Gopher")
	// ebiten.SetTPS(20) // Run Update() at 30hz instead of 60hz
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
