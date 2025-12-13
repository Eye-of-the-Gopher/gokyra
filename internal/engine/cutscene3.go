package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Scene 3 is the zooming map
type Scene3 struct {
	mapCounter int // 0 to 10 - map1, 10 - 20 - map2, 20 - 30 - map3, 30 - 40 - map4,40 - 50 - map5,50 - 60 - map6,
	map1       *ebiten.Image
	map2       *ebiten.Image
	map3       *ebiten.Image
	map4       *ebiten.Image
	map5       *ebiten.Image
	map6       *ebiten.Image
	pos        *ebiten.DrawImageOptions
}

func NewScene3(c *CutSceneManager) (*Scene3, error) {
	mapOrb1, err := c.assets.GetSprite("WTRDP1.CMP", "WTRDP2.COL", ScreenWidth, ScreenHeight, "")
	mapOrb2, err := c.assets.GetSprite("WTRDP2.CMP", "WTRDP2.COL", ScreenWidth, ScreenHeight, "")
	mapOrb3, err := c.assets.GetSprite("WTRDP3.CMP", "ORB.COL", ScreenWidth, ScreenHeight, "")

	if err != nil {
		EngineLogger.Error("Couldn't load sprite", "name", "WTRDP1.CMP")
		return nil, err
	}
	map1 := mapOrb1.GetImageRegion(0, 0, 160, 136).GetEbitenImage()
	map2 := mapOrb1.GetImageRegion(160, 0, 320, 136).GetEbitenImage()
	map3 := mapOrb2.GetImageRegion(0, 0, 160, 136).GetEbitenImage()
	map4 := mapOrb2.GetImageRegion(160, 0, 320, 136).GetEbitenImage()
	map5 := mapOrb3.GetImageRegion(0, 0, 160, 136).GetEbitenImage()
	map6 := mapOrb3.GetImageRegion(160, 0, 320, 136).GetEbitenImage()

	op := &ebiten.DrawImageOptions{}
	mapWidth := map1.Bounds().Dx()
	mapX := float64(ScreenWidth/2 - mapWidth/2)
	op.GeoM.Reset()
	op.GeoM.Translate(mapX, 10)

	return &Scene3{
		map1:       map1,
		map2:       map2,
		map3:       map3,
		map4:       map4,
		map5:       map5,
		map6:       map6,
		mapCounter: 0,
		pos:        op,
	}, nil
}

func (c *CutSceneManager) Scene3Update(game *Game) (bool, error) {
	if c.scene3.mapCounter < 60 {
		c.scene3.mapCounter += 1
	}
	return false, nil
}

func (c *CutSceneManager) Scene3Draw(screen *ebiten.Image, game *Game) {
	cntr := c.scene3.mapCounter
	var t *ebiten.Image
	if 0 <= cntr && cntr < 10 {
		t = c.scene3.map1
	}
	if 10 <= cntr && cntr < 20 {
		t = c.scene3.map2
	}
	if 20 <= cntr && cntr < 25 {
		t = c.scene3.map3
		c.subtitle = nil
	}
	if 25 <= cntr && cntr < 30 {
		t = c.scene3.map4
	}
	if 30 <= cntr && cntr < 35 {
		t = c.scene3.map5
	}
	if 35 <= cntr {
		t = c.scene3.map6
	}

	screen.DrawImage(t, c.scene3.pos)
}
