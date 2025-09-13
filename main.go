package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
	"github.com/nibrahim/eye-of-the-gopher/pkg/engine"
	"github.com/nibrahim/eye-of-the-gopher/pkg/formats"
)

func main() {
	formats.InitLogger(slog.LevelError, slog.LevelError, slog.LevelError, slog.LevelError)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage : %s [options] assetDirectory ...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nArguments:\n")
		fmt.Fprintf(os.Stderr, "  assetDirectory    Directory with original EOB .pak files (6 of them)\n")
	}
	extraAssetDir := flag.String("extraAssetDir", "", "Directory to side load extra assets")
	flag.Parse()

	if flag.NArg() == 0 {
		utils.ErrorAndExit("Error: No EOB origin asset directory specified")
	}

	assetDir := flag.Args()[0]
	game := engine.NewGame(assetDir, *extraAssetDir)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Eye Of The Gopher")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
