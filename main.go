package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
	assetLoader "github.com/nibrahim/eye-of-the-gopher/pkg/assets"
	"github.com/nibrahim/eye-of-the-gopher/pkg/formats"
)

type Game struct {
	assets formats.Assets
}

func (g *Game) assetDump() {
	g.assets.DumpAssets()
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func NewGame(assetDir string) Game {
	assets := assetLoader.LoadClassicAssets(assetDir)
	ret := Game{assets: *assets}
	return ret
}

func main() {
	utils.SetupLogging("eog.log")
	if len(os.Args) != 2 {
		utils.ErrorAndExit("Usage : %s asset_directory", os.Args[0])
	}

	assetDir := os.Args[1]
	game := NewGame(assetDir)
	eg := "DWARF.CPS"
	t, err := game.assets.GetSprite(eg, "")
	if err != nil {
		utils.ErrorAndExit("Couldn't load File %s: %v", eg, err)
	} else {
		fmt.Println(t)
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Eye Of The Gopher")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
