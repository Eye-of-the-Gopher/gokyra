package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/nibrahim/eye-of-the-gopher/pkg/assets"
	"github.com/nibrahim/eye-of-the-gopher/pkg/formats"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	assets formats.Assets
}

func NewGame() *Game {
	assets.LoadClassicAssets(assetDir string)
	ret := Game{}
	return &ret
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Eye of the Gopher\nHello, Dungeon!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	fmt.Println("Starting Eye of the Gopher...")

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Eye of the Gopher")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
