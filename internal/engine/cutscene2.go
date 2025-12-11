package engine

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
)

// Scene 2 is the orb
type Scene2 struct {
	orb        *ebiten.Image
	mageCircle *ebiten.Image
	mageFader  formats.FadeIterator
}

func NewScene2(c *CutSceneManager) (*Scene2, error) {
	towrmage, err := c.assets.GetSprite("TOWRMAGE.CMP", "TOWRMAGE.COL", 320, 200, "")
	if err != nil {
		EngineLogger.Error("Couldn't load sprite", "name", "towrmage.cmp")
		return nil, err
	}
	orb, err := c.assets.GetSprite("ORB.CMP", "ZOOMTUNL.COL", ScreenWidth, ScreenHeight, "")

	mageCircle := towrmage.GetImageRegion(128, 0, 256, 104)

	if err != nil {
		EngineLogger.Error("Couldn't load sprite", "name", "orb.cmp")
		return nil, err
	}
	return &Scene2{
		orb:       orb.GetImageRegion(0, 0, 160, 136).GetEbitenImage(),
		mageFader: mageCircle.GetEbitenImageFader(),
	}, nil
}

func (c *CutSceneManager) Scene2Update(game *Game) (bool, error) {
	c.subtitle = c.subtitles[2]
	if img, more := c.scene2.mageFader(); more {
		fmt.Println("Running here")
		c.scene2.mageCircle = img
	} else {
		fmt.Println("Running here - Done")
	}

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
