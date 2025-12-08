package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
)

// Actual scene 1 here
type Scene1 struct {
	towerSprite      image.Image
	windowSprite     image.Image
	mageCircleSprite image.Image
	textSprite       *formats.Sprite
	text1            image.Image
	scrollOffset     int
}

func NewScene1(c *CutSceneManager) (*Scene1, error) {
	textPalette, err := c.assets.GetPalette("TOWRMAGE.COL")
	if err != nil {
		EngineLogger.Error("Couldn't load palette for cutscene ", "palette", "TOWRMAGE.COL")
		return nil, err
	}
	textSprite, err := c.assets.GetSprite("TEXT.CMP", textPalette, 320, 200, c.scale, "")
	if err != nil {
		EngineLogger.Error("Couldn't load Text sprite", "sprite", "TEXT.CMP")
		return nil, err
	}

	srcRect := image.Rect(0, 0, 1280, 131)
	textSpriteImage, err := textSprite.GetEbitenImage()
	if err != nil {
		return nil, err
	}

	towrmage, err := c.assets.GetSprite("TOWRMAGE.CMP", textPalette, 320, 200, c.scale, "")
	towrmageImage, err := towrmage.GetEbitenImage()
	if err != nil {
		EngineLogger.Error("Couldn't load Towrmage sprite sheet", "sprite", "TOWRMAGE.CMP")
		return nil, err
	}
	// mageSprite, _ := c.assets.GetSprite("ENHANCED/X.PNG", textPalette, 120, 100, "")
	// i, _ := mageSprite.GetEbitenImage()

	// bounds := towerSprite.Image.Bounds()
	// fmt.Printf("Image dimensions: %dx%d\n", bounds.Dx(), bounds.Dy())

	// These rect dimentions were arrived at by
	// experimentation. 513 and 1023 seem to make sense. (512,
	// 1024) but not sure about the heights
	towerShaftRect := image.Rect(513, 417, 1023, 671)
	windowRect := image.Rect(0, 0, 510, 572)
	mageCircleRect := image.Rect(514, 0, 1023, 414)

	return &Scene1{
		towerSprite:      towrmageImage.SubImage(towerShaftRect),
		windowSprite:     towrmageImage.SubImage(windowRect),
		mageCircleSprite: towrmageImage.SubImage(mageCircleRect),
		textSprite:       textSprite,
		text1:            textSpriteImage.SubImage(srcRect),
		scrollOffset:     0,
	}, nil

}

func (c *CutSceneManager) Scene1Update(game *Game) (bool, error) {
	if c.subtitle == nil {
		c.subtitle = c.scene1.text1.(*ebiten.Image)
	}
	if c.scene1.scrollOffset <= 650 {
		c.scene1.scrollOffset += 1
	}

	return false, nil
}

func (c *CutSceneManager) Scene1Draw(screen *ebiten.Image, game *Game) {
	if c.scene1.mageCircleSprite != nil {
		op := &ebiten.DrawImageOptions{}
		t := c.scene1.mageCircleSprite.(*ebiten.Image)
		op.GeoM.Translate(357, 0)
		screen.DrawImage(t, op)
	}

	if c.scene1.towerSprite != nil {
		op := &ebiten.DrawImageOptions{}
		// op.GeoM.Translate(257, 200)
		t := c.scene1.towerSprite.(*ebiten.Image)
		op.GeoM.Translate(357, 415+float64(c.scene1.scrollOffset))
		screen.DrawImage(t, op)
		op.GeoM.Reset()
		op.GeoM.Translate(357, 163+float64(c.scene1.scrollOffset))
		screen.DrawImage(t, op)
		op.GeoM.Reset()
		op.GeoM.Translate(357, -89+float64(c.scene1.scrollOffset))
		screen.DrawImage(t, op)

	}

	if c.scene1.windowSprite != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(357, -661+float64(c.scene1.scrollOffset))
		t := c.scene1.windowSprite.(*ebiten.Image)
		screen.DrawImage(t, op)
	}
	if c.subtitle != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, 669)
		screen.DrawImage(c.subtitle, op)
	}

}
