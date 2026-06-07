package engine

import (
	"bytes"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
	"golang.org/x/image/font/gofont/goregular"
)

type charBox struct {
	x, y float32
}

var charSlots = [4]charBox{
	{16, 63},
	{70, 63},
}

type CharGenManager struct {
	bg      *ebiten.Image
	font    *text.GoTextFace
	pulse   float64
	current charBox
}

func NewCharGenManager(assets *formats.Assets) (*CharGenManager, error) {
	bg, err := assets.GetSprite("CHARGEN.CPS", "EOBPAL.COL", 320, 200, "")
	if err != nil {
		EngineLogger.Error("Couldn't load sprite", "name", "chargen.cps")
		return nil, fmt.Errorf("Couldn't load sprite chargen.cps")
	}
	source, _ := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	fontFace := &text.GoTextFace{
		Source: source,
		Size:   12,
	}

	ret := CharGenManager{
		current: charSlots[0],
		bg:      bg.GetEbitenImage(),
		font:    fontFace,
		pulse:   0.0,
	}
	return &ret, nil
}

func (c *CharGenManager) Update(game *Game) (GameState, error) {
	game.EnsureTrackPlaying("ENHANCED/CHARGEN.MP3")
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return GameQuit, nil
	}
	c.pulse += 0.1
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		c.current = charSlots[1]
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		c.current = charSlots[0]
	}

	return GameCharGen, nil

}

func (c *CharGenManager) Draw(screen *ebiten.Image, game *Game) {
	screen.DrawImage(c.bg, nil)
	textOp := &text.DrawOptions{}
	textOp.ColorScale.ScaleWithColor(color.White)
	textOp.GeoM.Translate(180, 100) // Starting position
	textOp.ColorScale.ScaleWithColor(color.White)
	textOp.LineSpacing = c.font.Metrics().HAscent + c.font.Metrics().HDescent
	text.Draw(screen, "Select the box of\nthe character you\nwish to create or\nview.", c.font, textOp)
	alpha := uint8((math.Sin(c.pulse) + 1) / 2 * 128)
	col := color.RGBA{0, 0, alpha / 2, alpha}
	// vector.StrokeRect(screen, c.current.x, c.current.y, 32, 33, 1, col, false)
	vector.DrawFilledRect(screen, c.current.x, c.current.y, 33, 34, col, false)

}
