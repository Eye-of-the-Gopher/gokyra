package engine

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
)

type CharGenManager struct {
	bg *ebiten.Image
}

func NewCharGenManager(assets *formats.Assets) (*CharGenManager, error) {
	bg, err := assets.GetSprite("CHARGEN.CPS", "EOBPAL.COL", 320, 200, "")
	if err != nil {
		EngineLogger.Error("Couldn't load sprite", "name", "chargen.cps")
		return nil, fmt.Errorf("Couldn't load sprite chargen.cps")
	}
	ret := CharGenManager{
		bg: bg.GetEbitenImage(),
	}
	return &ret, nil
}

func (c *CharGenManager) Update(game *Game) (GameState, error) {
	game.EnsureTrackPlaying("ENHANCED/CHARGEN.MP3")
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return GameQuit, nil
	}
	return GameCharGen, nil

}

func (c *CharGenManager) Draw(screen *ebiten.Image, game *Game) {
	screen.DrawImage(c.bg, nil)
}
