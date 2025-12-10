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
	towrmage, err := c.assets.GetSprite("TOWRMAGE.CMP", textPalette, 320, 200, "")
	if err != nil {
		EngineLogger.Error("Couldn't load Towrmage sprite sheet", "sprite", "TOWRMAGE.CMP")
		return nil, err
	}

	return &Scene1{
		towerSprite:      towrmage.GetEbitenImageRegion(128, 104, 256, 167),
		windowSprite:     towrmage.GetEbitenImageRegion(0, 0, 128, 143),
		mageCircleSprite: towrmage.GetEbitenImageRegion(128, 0, 256, 104),
		scrollOffset:     0,
	}, nil

}

func (c *CutSceneManager) Scene1Update(game *Game) (bool, error) {
	if c.subtitle == nil {
		c.subtitle = c.subtitles[0] // We the lord of waterdeep...
	}

	if c.scene1.scrollOffset <= 63*3+16 { // 63 is the height of the shaft segment. We move 3 segments down + 16 pixels for the extra in the window
		c.scene1.scrollOffset += 1
	}

	return false, nil
}

func (c *CutSceneManager) Scene1Draw(screen *ebiten.Image, game *Game) {
	if c.scene1.mageCircleSprite != nil {
		op := &ebiten.DrawImageOptions{}
		t := c.scene1.mageCircleSprite.(*ebiten.Image)
		op.GeoM.Translate(96, 30) // 16 pixels for the extra height of the window sprite and 14 for the black part on top of the mages
		screen.DrawImage(t, op)
	}

	if c.scene1.towerSprite != nil { // Tower shaft sprite height is 63. We have 4. 1 (top), 2, 3, and 4
		op := &ebiten.DrawImageOptions{}
		// op.GeoM.Translate(257, 200)
		t := c.scene1.towerSprite.(*ebiten.Image)
		op.GeoM.Translate(96, 126+float64(c.scene1.scrollOffset)) // Segment 4
		screen.DrawImage(t, op)
		op.GeoM.Reset()
		op.GeoM.Translate(96, 63+float64(c.scene1.scrollOffset)) // Segment 3
		screen.DrawImage(t, op)
		op.GeoM.Reset()
		op.GeoM.Translate(96, 0+float64(c.scene1.scrollOffset)) // Segment 2
		screen.DrawImage(t, op)
		op.GeoM.Reset()
		op.GeoM.Translate(96, -63+float64(c.scene1.scrollOffset)) // Segment 1
		screen.DrawImage(t, op)

	}

	if c.scene1.windowSprite != nil {
		op := &ebiten.DrawImageOptions{} // 143 is the height of the window. 63 for shaft segment 1
		op.GeoM.Translate(96, -143-63+float64(c.scene1.scrollOffset))
		t := c.scene1.windowSprite.(*ebiten.Image)
		screen.DrawImage(t, op)
	}
	if c.subtitle != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, 169)
		screen.DrawImage(c.subtitle, op)
	}

}
