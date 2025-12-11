package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Scene 2 is the orb
type Scene2 struct {
	orb        *ebiten.Image
	mageCircle *ebiten.Image
}

func NewScene2(c *CutSceneManager, mageCircle *ebiten.Image) (*Scene2, error) {
	orb, err := c.assets.GetSprite("ORB.CMP", "ZOOMTUNL.COL", ScreenWidth, ScreenHeight, "")
	if err != nil {
		EngineLogger.Error("Couldn't load sprite", "name", "orb.cmp")
		return nil, err
	}
	return &Scene2{
		orb:        orb.GetEbitenImageRegion(0, 0, 160, 136),
		mageCircle: mageCircle,
	}, nil
}

func (c *CutSceneManager) Scene2Update(game *Game) (bool, error) {
	c.subtitle = c.subtitles[2]
	return false, nil
}

func (c *CutSceneManager) Scene2Draw(screen *ebiten.Image, game *Game) {
	if c.scene1.mageCircleSprite != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(96, 20) // 16 pixels for the extra height of the window sprite and 14 for the black part on top of the mages
		screen.DrawImage(c.scene1.mageCircleSprite, op)
	}

	if c.scene2.orb != nil {
		op := &ebiten.DrawImageOptions{}
		orbWidth := c.scene2.orb.Bounds().Dx()
		orbX := float64(ScreenWidth/2 - orbWidth/2)
		t := c.scene2.orb
		op.GeoM.Translate(orbX, 10)
		screen.DrawImage(t, op)
	}

}
