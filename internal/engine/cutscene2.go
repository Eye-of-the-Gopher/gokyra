package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
)

// Scene 2 is the orb
type Scene2 struct {
	orb          *ebiten.Image
	mageCircle   *ebiten.Image
	map1         *ebiten.Image
	map1FaderIn  formats.FadeIterator
	mageFaderIn  formats.FadeIterator
	mageFaderOut formats.FadeIterator
	mageCounter  int
	orbCounter   int
	drawMap      bool
}

func NewScene2(c *CutSceneManager) (*Scene2, error) {
	towrmage, err := c.assets.GetSprite("TOWRMAGE.CMP", "TOWRMAGE.COL", 320, 200, "")
	if err != nil {
		EngineLogger.Error("Couldn't load sprite", "name", "towrmage.cmp")
		return nil, err
	}
	orb, err := c.assets.GetSprite("ORB.CMP", "ZOOMTUNL.COL", ScreenWidth, ScreenHeight, "")
	if err != nil {
		EngineLogger.Error("Couldn't load sprite", "name", "orb.cmp")
		return nil, err
	}
	mageCircle := towrmage.GetImageRegion(128, 0, 256, 104)

	mapOrb1, err := c.assets.GetSprite("WTRDP1.CMP", "WTRDP2.COL", ScreenWidth, ScreenHeight, "")
	if err != nil {
		EngineLogger.Error("Couldn't load sprite", "name", "WTRDP1.CMP")
		return nil, err
	}
	map1 := mapOrb1.GetImageRegion(0, 0, 160, 136)

	return &Scene2{
		orb:          orb.GetImageRegion(0, 0, 160, 136).GetEbitenImage(),
		map1FaderIn:  map1.GetEbitenImageFadeIn(10, 20),
		mageFaderIn:  mageCircle.GetEbitenImageFadeIn(10, 20),
		mageFaderOut: mageCircle.GetEbitenImageFadeOut(10, 20),
		mageCounter:  15,
		orbCounter:   10,
		drawMap:      false,
	}, nil
}

func (c *CutSceneManager) Scene2Update(game *Game) (bool, error) {
	c.subtitle = nil
	// Fade in the mages
	if img, mageFader := c.scene2.mageFaderIn(); mageFader {
		c.scene2.mageCircle = img
	} else {
		c.subtitle = c.subtitles[2] // Master!
		if c.scene2.mageCounter != 0 {
			c.scene2.mageCounter--
		}
	}

	// Fade out the mages
	if c.scene2.mageCounter == 0 {
		if img, mageFader := c.scene2.mageFaderOut(); mageFader {
			c.scene2.mageCircle = img
		} else {
			c.subtitle = c.subtitles[3] // They think they have found a solution
			if c.scene2.orbCounter != 0 {
				c.scene2.orbCounter--
				c.scene2.drawMap = true
			}
		}
	}
	if c.scene2.drawMap {
		c.scene2.mageCircle = nil
		c.scene2.orb = nil
		if img, map1Fader := c.scene2.map1FaderIn(); map1Fader {
			c.scene2.map1 = img
		} else {
			return true, nil
		}

	}

	return false, nil

}

func (c *CutSceneManager) Scene2Draw(screen *ebiten.Image, game *Game) {
	if c.scene2.mageCircle != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(96, 20)
		screen.DrawImage(c.scene2.mageCircle, op)
	}

	if c.scene2.orb != nil {
		op := &ebiten.DrawImageOptions{}
		orbWidth := c.scene2.orb.Bounds().Dx()
		orbX := float64(ScreenWidth/2 - orbWidth/2)
		t := c.scene2.orb
		op.GeoM.Translate(orbX, 10)
		screen.DrawImage(t, op)
	}

	if c.scene2.map1 != nil {
		op := &ebiten.DrawImageOptions{}
		t := c.scene2.map1
		orbWidth := t.Bounds().Dx()
		orbX := float64(ScreenWidth/2 - orbWidth/2)
		op.GeoM.Translate(orbX, 10)
		screen.DrawImage(t, op)
	}

}
