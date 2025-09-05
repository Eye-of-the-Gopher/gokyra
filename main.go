package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
	"github.com/nibrahim/eye-of-the-gopher/pkg/formats"
)

type Game struct {
	assets        formats.Assets
	image         *ebiten.Image
	audioContext  *audio.Context // ONE for entire game
	currentPlayer *audio.Player  // Changes per track

}

func (g *Game) assetDump() {
	g.assets.DumpAssets()
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.image, nil)
	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func NewGame(assetDir string, extraAssetDir string) Game {
	assets := formats.LoadAssets(assetDir, extraAssetDir)
	eg := "TITLE-V.CMP"
	pal := "WESTWOOD.COL"
	p, err := assets.GetPalette(pal)
	t, err := assets.GetSprite(eg, p, 320, 200, "")
	if err != nil {
		utils.ErrorAndExit("Couldn't load File %s: %v", eg, err)
	} else {
		fmt.Println(t)
	}
	audioContext, err := audio.NewContext(44100)
	if err != nil {
		return Game{
			assets:       *assets,
			image:        ebiten.NewImageFromImage(t.Image),
			audioContext: audioContext,
		}

	} else {
		return Game{
			assets:       *assets,
			image:        ebiten.NewImageFromImage(t.Image),
			audioContext: nil,
		}

	}
}

func main() {
	utils.SetupLogging("eog.log")
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
	game := NewGame(assetDir, *extraAssetDir)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Eye Of The Gopher")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
