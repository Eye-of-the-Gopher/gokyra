package engine

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
)

// Scene 5 is the adventurers in front of the king
type Scene5 struct {
	ground        *ebiten.Image
	groupArriving []*formats.Sprite
	groupIdx      int

	groundPos *ebiten.DrawImageOptions
	animIdx   int
}

func NewScene5(c *CutSceneManager) (*Scene5, error) {
	wtrdp3, err := c.assets.GetSprite("WTRDP3.CMP", "ORB.COL", ScreenWidth, ScreenHeight, "")

	if err != nil {
		EngineLogger.Error("Couldn't load sprite", "name", "WTRDP3.CMP")
		return nil, err
	}
	ground := wtrdp3.GetImageRegion(160, 0, 320, 136).GetEbitenImage()

	groupArriving := []*formats.Sprite{wtrdp3.GetImageRegion(0, 152, 35, 184),
		wtrdp3.GetImageRegion(41, 152, 76, 184),
		wtrdp3.GetImageRegion(81, 152, 116, 184),
		wtrdp3.GetImageRegion(41, 152, 76, 184),
	}

	op := &ebiten.DrawImageOptions{}
	mapWidth := ground.Bounds().Dx()
	mapX := float64(ScreenWidth/2 - mapWidth/2)
	op.GeoM.Reset()
	op.GeoM.Translate(mapX, 10)

	return &Scene5{
		ground:        ground,
		groupArriving: groupArriving,
		groundPos:     op,
		groupIdx:      0,
		animIdx:       0,
	}, nil
}

func (c *CutSceneManager) Scene5Update(game *Game) (bool, error) {
	c.scene4.animIdx += 1
	if c.scene4.animIdx < 75 {
		t := c.scene4.groupIdx
		t += 1
		t %= len(c.scene4.groupArriving)
		c.scene4.groupIdx = t

		fmt.Println("Ground counter ", t, c.scene4.groupX, c.scene4.groupY)
		c.scene4.groupX = c.scene4.groupX - 1
		c.scene4.groupY = c.scene4.groupY + 1
	}
	fmt.Println("Done")

	return false, nil
}

func (c *CutSceneManager) Scene5Draw(screen *ebiten.Image, game *Game) {
	if c.scene4.ground != nil {
		screen.DrawImage(c.scene4.ground, c.scene4.groundPos)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(c.scene4.groupX, c.scene4.groupY)
	grp := c.scene4.groupArriving[c.scene4.groupIdx].GetEbitenImage()
	screen.DrawImage(grp, op)

}
