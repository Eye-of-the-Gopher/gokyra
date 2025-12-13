package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
)

type PixelIterator func() (image.Point, bool)

// Framework here
type CutSceneManager struct {
	scene  int
	assets *formats.Assets

	subtitles []*ebiten.Image // All the subtitles pre loaded so that we can display them easily
	subtitle  *ebiten.Image   // The subtitle that is displayed. Cut out from the text sprite by each Update method and put here

	frameCntr  int
	frameDelay int // Will run update only these many times

	scene0 *Scene0
	scene1 *Scene1
	scene2 *Scene2
	scene3 *Scene3
}

func NewCutSceneManager(assets *formats.Assets) (*CutSceneManager, error) {
	subtitles, err := loadSubtitles(assets)
	if err != nil {
		return nil, err
	}

	csm := &CutSceneManager{scene: 0,
		assets:     assets,
		frameCntr:  0,
		frameDelay: 5,
		subtitles:  subtitles,
	}

	sm0, err := NewScene0(csm)
	sm1, err := NewScene1(csm)
	sm2, err := NewScene2(csm)
	sm3, err := NewScene3(csm)

	if err != nil {
		return nil, err
	}
	csm.scene0 = sm0
	csm.scene1 = sm1
	csm.scene2 = sm2
	csm.scene3 = sm3
	return csm, nil

}

func (c *CutSceneManager) Update(game *Game) (bool, error) {
	c.frameCntr += 1
	if c.frameCntr == c.frameDelay {
		switch c.scene {
		case 0:
			next, _ := c.Scene0Update(game)
			if next {
				c.scene = 1
			}
		case 1:
			next, _ := c.Scene1Update(game)
			if next {
				c.scene = 2 // Go to next
			}
		case 2:
			next, _ := c.Scene2Update(game)
			if next {
				c.scene = 3 // Go to next
			}
		case 3:
			next, _ := c.Scene3Update(game)
			if next {
				c.scene = 3 // Go to next
			}
		default:
			EngineLogger.Warn("Scene not implemented yet", "scene", c.scene)
		}
		c.frameCntr = 0
	}
	return false, nil

}

func (c *CutSceneManager) Draw(screen *ebiten.Image, game *Game) {
	switch c.scene {
	case 0:
		c.Scene0Draw(screen, game)
	case 1:
		c.Scene1Draw(screen, game)
	case 2:
		c.Scene2Draw(screen, game)
	case 3:
		c.Scene3Draw(screen, game)
	default:
		EngineLogger.Warn("Scene not implemented yet", "scene", c.scene)
	}
	if c.subtitle != nil {
		op := &ebiten.DrawImageOptions{}
		subtitleHeight := c.subtitle.Bounds().Dy()
		screenHeight := screen.Bounds().Dy()
		op.GeoM.Translate(0, float64(screenHeight-subtitleHeight))
		screen.DrawImage(c.subtitle, op)
	}

}

// Helpers

func loadSubtitles(assets *formats.Assets) ([]*ebiten.Image, error) {
	subtitleSprite, err := assets.GetSprite("TEXT.CMP", "TOWRMAGE.COL", 320, 200, "")
	if err != nil {
		EngineLogger.Error("Couldn't load Text sprite", "sprite", "TEXT.CMP")
		return nil, err
	}
	ret := []*ebiten.Image{}

	ret = append(ret, subtitleSprite.GetImageRegion(0, 0, 320, 32).GetEbitenImage())    // 0 We the lord of waterdeep...
	ret = append(ret, subtitleSprite.GetImageRegion(0, 32, 320, 67).GetEbitenImage())   // 1 Give call to the heroes of the land...
	ret = append(ret, subtitleSprite.GetImageRegion(0, 67, 320, 80).GetEbitenImage())   // 2 Master!
	ret = append(ret, subtitleSprite.GetImageRegion(0, 80, 320, 96).GetEbitenImage())   // 3 They think they have found...
	ret = append(ret, subtitleSprite.GetImageRegion(0, 96, 320, 128).GetEbitenImage())  // 4 We commission you to find...
	ret = append(ret, subtitleSprite.GetImageRegion(0, 128, 320, 144).GetEbitenImage()) // 5 Prepare for the ...
	ret = append(ret, subtitleSprite.GetImageRegion(0, 144, 320, 160).GetEbitenImage()) // 6 Begin your search below ...
	ret = append(ret, subtitleSprite.GetImageRegion(0, 160, 320, 175).GetEbitenImage()) // 7 We have them...
	ret = append(ret, subtitleSprite.GetImageRegion(0, 175, 320, 191).GetEbitenImage()) // 8 Their fate is sealed...

	return ret, nil
}

func fadeGridGen(x0, y0, edge int) PixelIterator {
	i := 0
	dir := 0
	pt := image.Point{X: x0, Y: y0}

	return func() (image.Point, bool) {
		if edge <= 1 {
			EngineLogger.Debug("We're done. Stopping now")
			return image.Point{}, false
		}
		// EngineLogger.Debug("PixelIterator ", "i", i, "dir", dir, "point", pt)
		if i >= edge {
			dir += 1
			i = 0
			// EngineLogger.Debug("Resetting i to 0", "dir", dir)
		}

		if dir == 4 {
			edge -= 1
			i = 0
			dir = 0
			pt = image.Point{X: x0, Y: y0}
			// EngineLogger.Debug("Reducing edge", "edge", edge)
		}

		switch dir {
		case 0:
			// EngineLogger.Debug("Incrementing X")
			pt.X += 1
		case 1:
			// EngineLogger.Debug("Incrementing Y")
			pt.Y += 1
		case 2:
			// EngineLogger.Debug("Decrementing X")
			pt.X -= 1
		case 3:
			// EngineLogger.Debug("Decrementing Y")
			pt.Y -= 1
		}
		i++
		return pt, true // (value, has_more)

	}

}
