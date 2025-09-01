package main

import (
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
	"github.com/nibrahim/eye-of-the-gopher/pkg/formats"
)

type Game struct {
	assets formats.Assets
	image  *ebiten.Image
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

func NewGame(assetDir string) Game {
	assets := formats.LoadClassicAssets(assetDir)
	eg := "TITLE-V.CMP"
	pal := "WESTWOOD.COL"
	p, err := assets.GetPalette(pal)
	t, err := assets.GetSprite(eg, p, 320, 200, "")
	if err != nil {
		utils.ErrorAndExit("Couldn't load File %s: %v", eg, err)
	} else {
		fmt.Println(t)
	}

	f, _ := os.Create("/tmp/baz.png")
	defer f.Close()
	png.Encode(f, t.Image)

	ret := Game{
		assets: *assets,
		image:  ebiten.NewImageFromImage(t.Image),
	}

	return ret
}

func main() {
	utils.SetupLogging("eog.log")
	if len(os.Args) != 2 {
		utils.ErrorAndExit("Usage : %s asset_directory", os.Args[0])
	}

	assetDir := os.Args[1]
	game := NewGame(assetDir)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Eye Of The Gopher")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
