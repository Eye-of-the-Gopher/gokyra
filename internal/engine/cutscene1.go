package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
)

// Actual scene 1 here
type Scene1 struct {
	towerSprite      *ebiten.Image
	windowSprite     *ebiten.Image
	mageCircleSprite *ebiten.Image
	textSprite       *formats.Sprite
	text1            image.Image
	scrollOffset     int
	holdCounter      int
	holdLimit        int
}

func NewScene1(c *CutSceneManager) (*Scene1, error) {
	towrmage, err := c.assets.GetSprite("TOWRMAGE.CMP", "TOWRMAGE.COL", 320, 200, "")
	if err != nil {
		EngineLogger.Error("Couldn't load Towrmage sprite sheet", "sprite", "TOWRMAGE.CMP")
		return nil, err
	}

	return &Scene1{
		towerSprite:      towrmage.GetImageRegion(128, 104, 256, 167).GetEbitenImage(),
		windowSprite:     towrmage.GetImageRegion(0, 0, 128, 143).GetEbitenImage(),
		mageCircleSprite: towrmage.GetImageRegion(128, 0, 256, 104).GetEbitenImage(),
		scrollOffset:     0,
		holdCounter:      0,
		holdLimit:        15,
	}, nil

}

func (c *CutSceneManager) Scene1Update(game *Game) (bool, error) {
	if c.subtitle == nil {
		c.subtitle = c.subtitles[0] // We the lord of waterdeep...
	}

	if c.scene1.scrollOffset <= 63*2+16 { // 63 is the height of the shaft segment. We move 2 segments down + 16 pixels for the extra in the window
		c.scene1.scrollOffset += 1
	} else {
		c.scene1.holdCounter += 1
		c.subtitle = c.subtitles[1] // Give call to the heroes of the land...
	}

	if c.scene1.holdCounter > c.scene1.holdLimit {
		return true, nil
	}
	return false, nil

}

func (c *CutSceneManager) Scene1Draw(screen *ebiten.Image, game *Game) {
	if c.scene1.mageCircleSprite != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(96, 20) // 16 pixels for the extra height of the window sprite and 14 for the black part on top of the mages
		screen.DrawImage(c.scene1.mageCircleSprite, op)
	}

	if c.scene1.towerSprite != nil { // Tower shaft sprite height is 63. We have 4. 2 (top), 3, and 4
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(96, 126+float64(c.scene1.scrollOffset)) // Segment 4
		screen.DrawImage(c.scene1.towerSprite, op)
		op.GeoM.Reset()
		op.GeoM.Translate(96, 63+float64(c.scene1.scrollOffset)) // Segment 3
		screen.DrawImage(c.scene1.towerSprite, op)
		op.GeoM.Reset()
		op.GeoM.Translate(96, 0+float64(c.scene1.scrollOffset)) // Segment 2
		screen.DrawImage(c.scene1.towerSprite, op)

	}

	if c.scene1.windowSprite != nil {
		op := &ebiten.DrawImageOptions{} // 143 is the height of the window. 63 for shaft segment 1
		op.GeoM.Translate(96, -143+float64(c.scene1.scrollOffset))
		screen.DrawImage(c.scene1.windowSprite, op)
	}

}
