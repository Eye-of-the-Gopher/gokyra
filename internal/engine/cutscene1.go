package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
)

// Actual scene 1 here
type Scene1 struct {
	towerSprite *formats.Sprite
	textSprite  *formats.Sprite
	text1       image.Image
}

func NewScene1(c *CutSceneManager) (*Scene1, error) {
	textPalette, err := c.assets.GetPalette("TOWRMAGE.COL")
	if err != nil {
		EngineLogger.Error("Couldn't load palette for cutscene ", "palette", "TOWRMAGE.COL")
		return nil, err
	}
	textSprite, err := c.assets.GetSprite("TEXT.CMP", textPalette, 320, 200, "")
	if err != nil {
		EngineLogger.Error("Couldn't load Text sprite", "sprite", "TEXT.CMP")
		return nil, err
	}

	srcRect := image.Rect(0, 0, 1280, 131)
	textSpriteImage, err := textSprite.GetEbitenImage()
	if err != nil {
		return nil, err
	}

	return &Scene1{
		towerSprite: nil,
		textSprite:  textSprite,
		text1:       textSpriteImage.SubImage(srcRect).(*ebiten.Image),
	}, nil

}

func (c *CutSceneManager) Scene1Update(game *Game) (bool, error) {
	if c.subtitle == nil {
		c.subtitle = c.scene1.text1.(*ebiten.Image)
	}

	return false, nil
}

func (c *CutSceneManager) Scene1Draw(screen *ebiten.Image, game *Game) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 669)
	if c.subtitle != nil {
		screen.DrawImage(c.subtitle, op)

	}
}
