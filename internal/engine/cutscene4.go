package engine

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
)

// Scene 4 is the people coming up to the king
type Scene4 struct {
	groupX        float64
	groupY        float64
	ground        *ebiten.Image
	groupArriving []*formats.Sprite
	groupLeaving  []*formats.Sprite
	groundPos     *ebiten.DrawImageOptions
	groundCntr    int
	tmp           bool
}

func NewScene4(c *CutSceneManager) (*Scene4, error) {
	wtrdp3, err := c.assets.GetSprite("WTRDP3.CMP", "ORB.COL", ScreenWidth, ScreenHeight, "")

	if err != nil {
		EngineLogger.Error("Couldn't load sprite", "name", "WTRDP3.CMP")
		return nil, err
	}
	ground := wtrdp3.GetImageRegion(160, 0, 320, 136).GetEbitenImage()

	groupArriving := []*formats.Sprite{wtrdp3.GetImageRegion(0, 152, 35, 184),
		wtrdp3.GetImageRegion(41, 152, 76, 184),
		wtrdp3.GetImageRegion(81, 152, 116, 184),
	}

	op := &ebiten.DrawImageOptions{}
	mapWidth := ground.Bounds().Dx()
	mapX := float64(ScreenWidth/2 - mapWidth/2)
	op.GeoM.Reset()
	op.GeoM.Translate(mapX, 10)

	return &Scene4{
		groupX:        185,
		groupY:        17,
		ground:        ground,
		groupArriving: groupArriving,
		groundPos:     op,
		groundCntr:    0,
	}, nil
}

func (c *CutSceneManager) Scene4Update(game *Game) (bool, error) {
	if c.scene4.groupX > 110 {
		if c.scene4.tmp {
			c.scene4.groundCntr = c.scene4.groundCntr + 1
			c.scene4.tmp = false
		} else {
			c.scene4.tmp = true
		}
		fmt.Println("Ground counter ", c.scene4.groundCntr, c.scene4.groupX, c.scene4.groupY)
		c.scene4.groundCntr = c.scene4.groundCntr % 3
		c.scene4.groupX = c.scene4.groupX - 1
		c.scene4.groupY = c.scene4.groupY + 1
	}

	return false, nil
}

func (c *CutSceneManager) Scene4Draw(screen *ebiten.Image, game *Game) {
	if c.scene4.ground != nil {
		screen.DrawImage(c.scene4.ground, c.scene4.groundPos)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(c.scene4.groupX, c.scene4.groupY)
	grp := c.scene4.groupArriving[c.scene4.groundCntr].GetEbitenImage()
	screen.DrawImage(grp, op)

}
