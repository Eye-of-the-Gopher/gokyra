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

	if err != nil {
		return nil, err
	}
	csm.scene0 = sm0
	csm.scene1 = sm1
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
				c.scene = 1 // Go to next
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
	textPalette, err := assets.GetPalette("TOWRMAGE.COL")
	if err != nil {
		EngineLogger.Error("Couldn't load palette for cutscene ", "palette", "TOWRMAGE.COL")
		return nil, err
	}

	subtitleSprite, err := assets.GetSprite("TEXT.CMP", textPalette, 320, 200, "")
	if err != nil {
		EngineLogger.Error("Couldn't load Text sprite", "sprite", "TEXT.CMP")
		return nil, err
	}
	ret := []*ebiten.Image{}
	ret = append(ret, subtitleSprite.GetEbitenImageRegion(0, 0, 320, 31))    // We the lord of waterdeep...
	ret = append(ret, subtitleSprite.GetEbitenImageRegion(0, 32, 320, 67))   // Give call to the heroes of the land...
	ret = append(ret, subtitleSprite.GetEbitenImageRegion(0, 67, 320, 80))   // Master..
	ret = append(ret, subtitleSprite.GetEbitenImageRegion(0, 80, 320, 96))   // They think they have found...
	ret = append(ret, subtitleSprite.GetEbitenImageRegion(0, 96, 320, 128))  // We commission you to find...
	ret = append(ret, subtitleSprite.GetEbitenImageRegion(0, 128, 320, 144)) // Prepare for the ...
	ret = append(ret, subtitleSprite.GetEbitenImageRegion(0, 144, 320, 160)) // Begin your search below ...
	ret = append(ret, subtitleSprite.GetEbitenImageRegion(0, 160, 320, 175)) // We have them...
	ret = append(ret, subtitleSprite.GetEbitenImageRegion(0, 175, 320, 191)) // Their fate is sealed...

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
