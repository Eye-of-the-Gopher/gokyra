package engine

import (
	"bytes"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
)

type MenuManager struct {
	background     *ebiten.Image
	selectedOption int
	menuItems      [3]string
	menubox        *ebiten.Image
	font           *text.GoTextFace
}

func NewMenuManager(assets *formats.Assets) *MenuManager {
	EngineLogger.Debug("Creating Menu Manager")
	background, err := assets.GetSprite("INTRO.CPS", "EOBPAL.COL", 320, 200, "")
	if err != nil {
		EngineLogger.Error("Couldn't load  background sprite", "sprite", "INTRO.CPS")
		return nil
	}

	menuItems := [3]string{"Load game in progress", "Start a new party", "Exit"}
	menuRect := ebiten.NewImage(240, 35)
	menuRect.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 180})

	source, _ := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	fontFace := &text.GoTextFace{
		Source: source,
		Size:   12,
	}

	mm := MenuManager{
		background:     background.GetEbitenImage(),
		menuItems:      menuItems,
		selectedOption: 0,
		menubox:        menuRect,
		font:           fontFace,
	}
	return &mm

}

func (m *MenuManager) Update(g *Game) (bool, error) {
	// Key just pressed this frame (not held)
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		m.selectedOption++
		if m.selectedOption >= len(m.menuItems) {
			m.selectedOption = 0 // Wrap around
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		m.selectedOption--
		if m.selectedOption < 0 {
			m.selectedOption = len(m.menuItems) - 1 // Wrap around
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		EngineLogger.Debug("Selected item ", "entry", m.menuItems[m.selectedOption])
		if m.selectedOption == 2 {
			return true, nil
		}
	}

	return false, nil
}

func (m *MenuManager) Draw(screen *ebiten.Image, g *Game) (bool, error) {
	screen.DrawImage(m.background, nil)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(40, 162)
	screen.DrawImage(m.menubox, op)

	textOp := &text.DrawOptions{}
	textOp.ColorScale.ScaleWithColor(color.White)
	textOp.GeoM.Translate(50, 162) // Starting position
	for idx, menuItem := range m.menuItems {
		textOp.ColorScale.Reset()
		if idx == m.selectedOption {
			textOp.ColorScale.ScaleWithColor(color.RGBA{R: 200, G: 0, B: 0, A: 255})
		} else {
			textOp.ColorScale.ScaleWithColor(color.White)
		}
		text.Draw(screen, menuItem, m.font, textOp)
		textOp.GeoM.Translate(0, 10)
	}

	return false, nil

}
